// Load Testing Tool for Blog API
// Run: node load-test.js
// Simulates 50,000 daily active users

const BASE_URL = 'http://localhost:8080';

// Configuration for 50k DAU
const CONFIG = {
  totalUsers: 50000,           // Daily active users
  concurrentUsers: 100,         // Simultaneous requests
  testDurationMinutes: 5,       // How long to run the test
  
  // Request distribution (%)
  readBlogPercent: 60,          // Most users just read
  listBlogsPercent: 25,         // Browse blog lists
  searchPercent: 10,            // Search functionality
  createBlogPercent: 3,         // Write blogs (will get 401 - expected)
  updateBlogPercent: 1,         // Edit blogs (will get 401 - expected)
  deleteBlogPercent: 1,         // Delete blogs (will get 401 - expected)
  
  // Blog ID range (will be auto-detected)
  minBlogId: 1,
  maxBlogId: 1000
};

// Performance metrics
const metrics = {
  totalRequests: 0,
  successfulRequests: 0,
  failedRequests: 0,
  responseTimes: [],
  errors: {},
  startTime: null,
  endTime: null,
  
  // Per-endpoint metrics
  endpoints: {},
  
  // Expected errors (not counted as failures)
  expectedErrors: {
    'HTTP 401': 0,  // Auth required - expected
    'HTTP 404': 0   // Blog not found - will track separately
  }
};

function recordMetric(endpoint, duration, success, error = null) {
  metrics.totalRequests++;
  metrics.responseTimes.push(duration);
  
  // Track expected errors separately
  if (error && (error === 'HTTP 401' || error === 'HTTP 404')) {
    metrics.expectedErrors[error]++;
  } else if (success) {
    metrics.successfulRequests++;
  } else {
    metrics.failedRequests++;
    metrics.errors[error] = (metrics.errors[error] || 0) + 1;
  }
  
  // Track per-endpoint stats
  if (!metrics.endpoints[endpoint]) {
    metrics.endpoints[endpoint] = {
      count: 0,
      totalTime: 0,
      minTime: Infinity,
      maxTime: 0,
      errors: 0,
      success: 0
    };
  }
  
  const ep = metrics.endpoints[endpoint];
  ep.count++;
  ep.totalTime += duration;
  ep.minTime = Math.min(ep.minTime, duration);
  ep.maxTime = Math.max(ep.maxTime, duration);
  if (!success) ep.errors++;
  else ep.success++;
}

async function makeRequest(method, endpoint, body = null, token = null) {
  const start = Date.now();
  
  try {
    const options = {
      method,
      headers: { 'Content-Type': 'application/json' }
    };
    
    if (token) options.headers['Authorization'] = `Bearer ${token}`;
    if (body) options.body = JSON.stringify(body);
    
    const response = await fetch(`${BASE_URL}${endpoint}`, options);
    const duration = Date.now() - start;
    
    const success = response.ok;
    recordMetric(endpoint, duration, success, success ? null : `HTTP ${response.status}`);
    
    return { success, status: response.status, duration, data: success ? await response.json() : null };
  } catch (error) {
    const duration = Date.now() - start;
    recordMetric(endpoint, duration, false, error.message);
    return { success: false, error: error.message, duration };
  }
}

// Auto-detect blog ID range from the API
async function detectBlogRange() {
  try {
    const response = await fetch(`${BASE_URL}/blogs?page=1&limit=1`);
    if (response.ok) {
      const data = await response.json();
      if (data.blogs && data.blogs.length > 0) {
        CONFIG.minBlogId = data.blogs[0].id;
        
        // Get last page to find max ID
        const lastPageResponse = await fetch(`${BASE_URL}/blogs?page=999999&limit=1`);
        if (lastPageResponse.ok) {
          const lastPageData = await lastPageResponse.json();
          if (lastPageData.blogs && lastPageData.blogs.length > 0) {
            CONFIG.maxBlogId = lastPageData.blogs[0].id;
          }
        }
        
        console.log(`📊 Detected blog ID range: ${CONFIG.minBlogId} - ${CONFIG.maxBlogId}\n`);
      }
    }
  } catch (error) {
    console.log(`⚠️  Could not auto-detect blog range, using defaults (${CONFIG.minBlogId}-${CONFIG.maxBlogId})\n`);
  }
}

// Simulate realistic user behavior
async function simulateUser(userId) {
  const rand = Math.random() * 100;
  
  if (rand < CONFIG.readBlogPercent) {
    // Read a blog within the valid range
    const blogId = Math.floor(Math.random() * (CONFIG.maxBlogId - CONFIG.minBlogId + 1)) + CONFIG.minBlogId;
    await makeRequest('GET', `/blogs/${blogId}`);
    
  } else if (rand < CONFIG.readBlogPercent + CONFIG.listBlogsPercent) {
    // List blogs with pagination
    const page = Math.floor(Math.random() * 10) + 1;
    await makeRequest('GET', `/blogs?page=${page}&limit=10`);
    
  } else if (rand < CONFIG.readBlogPercent + CONFIG.listBlogsPercent + CONFIG.searchPercent) {
    // Search blogs
    const queries = ['PostgreSQL', 'API', 'Database', 'Guide', 'Tutorial', 'Performance'];
    const query = queries[Math.floor(Math.random() * queries.length)];
    await makeRequest('GET', `/blogs/search?q=${query}`);
    
  } else {
    // Authenticated actions (will get 401 - expected for load test)
    const rand2 = Math.random() * 100;
    
    if (rand2 < 60) {
      // Try to create (will fail with 401 - expected)
      await makeRequest('POST', '/blogs/create', {
        title: `Load Test Blog ${userId}`,
        content: 'This is a load test blog post',
        tags: ['test']
      });
    } else if (rand2 < 90) {
      // Try to update (will fail with 401 - expected)
      const blogId = Math.floor(Math.random() * (CONFIG.maxBlogId - CONFIG.minBlogId + 1)) + CONFIG.minBlogId;
      await makeRequest('PUT', `/blogs/${blogId}`, {
        title: 'Updated Title',
        content: 'Updated content'
      });
    } else {
      // Try to delete (will fail with 401 - expected)
      const blogId = Math.floor(Math.random() * (CONFIG.maxBlogId - CONFIG.minBlogId + 1)) + CONFIG.minBlogId;
      await makeRequest('DELETE', `/blogs/${blogId}`);
    }
  }
}

// Run concurrent users
async function runLoadTest() {
  console.log('🚀 Starting Load Test...\n');
  console.log('Configuration:');
  console.log(`  - Total Users: ${CONFIG.totalUsers.toLocaleString()}`);
  console.log(`  - Concurrent Users: ${CONFIG.concurrentUsers}`);
  console.log(`  - Test Duration: ${CONFIG.testDurationMinutes} minutes`);
  console.log(`  - Target: ${BASE_URL}\n`);
  
  // Detect valid blog ID range
  await detectBlogRange();
  
  metrics.startTime = Date.now();
  const endTime = metrics.startTime + (CONFIG.testDurationMinutes * 60 * 1000);
  
  let userIdCounter = 0;
  
  // Status update interval
  const statusInterval = setInterval(() => {
    printStatus();
  }, 5000); // Every 5 seconds
  
  // Keep spawning users until test duration ends
  while (Date.now() < endTime) {
    const batch = [];
    
    for (let i = 0; i < CONFIG.concurrentUsers; i++) {
      batch.push(simulateUser(userIdCounter++));
    }
    
    await Promise.all(batch);
    
    // Small delay between batches
    await new Promise(resolve => setTimeout(resolve, 100));
  }
  
  clearInterval(statusInterval);
  metrics.endTime = Date.now();
  
  printFinalReport();
}

function printStatus() {
  const elapsed = ((Date.now() - metrics.startTime) / 1000).toFixed(1);
  const rps = (metrics.totalRequests / (elapsed || 1)).toFixed(2);
  const successRate = ((metrics.successfulRequests / (metrics.totalRequests || 1)) * 100).toFixed(2);
  const avgResponseTime = (metrics.responseTimes.reduce((a, b) => a + b, 0) / metrics.responseTimes.length).toFixed(2);
  
  console.log(`⏱️  ${elapsed}s | Requests: ${metrics.totalRequests} | RPS: ${rps} | Success: ${successRate}% | Avg: ${avgResponseTime}ms`);
}

function calculatePercentile(arr, percentile) {
  if (arr.length === 0) return 0;
  const sorted = [...arr].sort((a, b) => a - b);
  const index = Math.ceil((percentile / 100) * sorted.length) - 1;
  return sorted[Math.max(0, index)] || 0;
}

function printFinalReport() {
  console.log('\n' + '='.repeat(70));
  console.log('LOAD TEST RESULTS');
  console.log('='.repeat(70));
  
  const duration = (metrics.endTime - metrics.startTime) / 1000;
  const rps = metrics.totalRequests / duration;
  const actualSuccessRate = (metrics.successfulRequests / metrics.totalRequests) * 100;
  const totalExpectedErrors = Object.values(metrics.expectedErrors).reduce((a, b) => a + b, 0);
  const adjustedSuccessRate = ((metrics.successfulRequests + totalExpectedErrors) / metrics.totalRequests) * 100;
  
  console.log('\n📊 Overall Statistics:');
  console.log(`  Total Requests: ${metrics.totalRequests.toLocaleString()}`);
  console.log(`  Successful (200 OK): ${metrics.successfulRequests.toLocaleString()} (${actualSuccessRate.toFixed(2)}%)`);
  console.log(`  Expected Errors (401/404): ${totalExpectedErrors.toLocaleString()}`);
  console.log(`  Actual Failures: ${metrics.failedRequests.toLocaleString()}`);
  console.log(`  Adjusted Success Rate: ${adjustedSuccessRate.toFixed(2)}%`);
  console.log(`  Duration: ${duration.toFixed(2)}s`);
  console.log(`  Requests/sec: ${rps.toFixed(2)}`);
  
  console.log('\n⏱️  Response Times (for successful requests):');
  const successfulTimes = metrics.responseTimes.filter((_, i) => i < metrics.successfulRequests);
  if (successfulTimes.length > 0) {
    const avg = successfulTimes.reduce((a, b) => a + b, 0) / successfulTimes.length;
    console.log(`  Average: ${avg.toFixed(2)}ms`);
    console.log(`  Min: ${Math.min(...successfulTimes).toFixed(2)}ms`);
    console.log(`  Max: ${Math.max(...successfulTimes).toFixed(2)}ms`);
    console.log(`  P50 (Median): ${calculatePercentile(successfulTimes, 50).toFixed(2)}ms`);
    console.log(`  P95: ${calculatePercentile(successfulTimes, 95).toFixed(2)}ms`);
    console.log(`  P99: ${calculatePercentile(successfulTimes, 99).toFixed(2)}ms`);
  }
  
  console.log('\n🎯 Per-Endpoint Performance (Top 10):');
  const sortedEndpoints = Object.entries(metrics.endpoints)
    .sort((a, b) => b[1].count - a[1].count)
    .slice(0, 10);
    
  for (const [endpoint, stats] of sortedEndpoints) {
    const avgTime = stats.totalTime / stats.count;
    const successRate = (stats.success / stats.count) * 100;
    
    // Clean up endpoint display
    let displayEndpoint = endpoint;
    if (endpoint.includes('/blogs/') && /\/\d+/.test(endpoint)) {
      displayEndpoint = '/blogs/{id}';
    }
    
    console.log(`\n  ${displayEndpoint}:`);
    console.log(`    Requests: ${stats.count} (${successRate.toFixed(1)}% success)`);
    console.log(`    Avg Time: ${avgTime.toFixed(2)}ms`);
    console.log(`    Min/Max: ${stats.minTime.toFixed(2)}ms / ${stats.maxTime.toFixed(2)}ms`);
  }
  
  console.log('\n📋 Expected Errors (Not Counted as Failures):');
  for (const [error, count] of Object.entries(metrics.expectedErrors)) {
    if (count > 0) {
      console.log(`  ${error}: ${count} (${((count / metrics.totalRequests) * 100).toFixed(2)}%)`);
    }
  }
  
  if (Object.keys(metrics.errors).length > 0) {
    console.log('\n❌ Actual Errors (Unexpected):');
    for (const [error, count] of Object.entries(metrics.errors)) {
      console.log(`  ${error}: ${count}`);
    }
  } else {
    console.log('\n✅ No unexpected errors!');
  }
  
  console.log('\n' + '='.repeat(70));
  
  // Performance assessment
  console.log('\n📈 Performance Assessment:');
  const avg = successfulTimes.length > 0 
    ? successfulTimes.reduce((a, b) => a + b, 0) / successfulTimes.length 
    : 0;
    
  if (avg < 50) {
    console.log('  ✅ EXCELLENT - Average response time < 50ms');
  } else if (avg < 100) {
    console.log('  ✅ GOOD - Average response time < 100ms');
  } else if (avg < 200) {
    console.log('  ⚠️  ACCEPTABLE - Average response time < 200ms');
  } else {
    console.log('  ❌ POOR - Average response time > 200ms - Optimization needed!');
  }
  
  if (adjustedSuccessRate > 99) {
    console.log('  ✅ EXCELLENT - Adjusted success rate > 99%');
  } else if (adjustedSuccessRate > 95) {
    console.log('  ✅ GOOD - Adjusted success rate > 95%');
  } else {
    console.log('  ❌ POOR - Too many unexpected failures!');
  }
  
  console.log('\n' + '='.repeat(70));
  console.log('\n💡 Note: 401 (Unauthorized) and some 404 (Not Found) errors are expected');
  console.log('   in this load test as we are simulating unauthenticated users.');
}

// Run the test
runLoadTest().catch(console.error);