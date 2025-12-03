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
  createBlogPercent: 3,         // Write blogs
  updateBlogPercent: 1,         // Edit blogs
  deleteBlogPercent: 1,         // Delete blogs
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
  endpoints: {}
};

function recordMetric(endpoint, duration, success, error = null) {
  metrics.totalRequests++;
  metrics.responseTimes.push(duration);
  
  if (success) {
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
      errors: 0
    };
  }
  
  const ep = metrics.endpoints[endpoint];
  ep.count++;
  ep.totalTime += duration;
  ep.minTime = Math.min(ep.minTime, duration);
  ep.maxTime = Math.max(ep.maxTime, duration);
  if (!success) ep.errors++;
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
    
    return { success, status: response.status, duration };
  } catch (error) {
    const duration = Date.now() - start;
    recordMetric(endpoint, duration, false, error.message);
    return { success: false, error: error.message, duration };
  }
}

// Simulate realistic user behavior
async function simulateUser(userId) {
  const rand = Math.random() * 100;
  
  if (rand < CONFIG.readBlogPercent) {
    // Read a random blog
    const blogId = Math.floor(Math.random() * 1000) + 1;
    await makeRequest('GET', `/blogs/${blogId}`);
    
  } else if (rand < CONFIG.readBlogPercent + CONFIG.listBlogsPercent) {
    // List blogs with pagination
    const page = Math.floor(Math.random() * 10) + 1;
    await makeRequest('GET', `/blogs?page=${page}&limit=10`);
    
  } else if (rand < CONFIG.readBlogPercent + CONFIG.listBlogsPercent + CONFIG.searchPercent) {
    // Search blogs
    const queries = ['tech', 'tutorial', 'news', 'guide', 'review'];
    const query = queries[Math.floor(Math.random() * queries.length)];
    await makeRequest('GET', `/blogs/search?q=${query}`);
    
  } else {
    // Authenticated actions (create/update/delete)
    // For simplicity, we'll just simulate the request pattern
    // In real scenario, you'd need valid tokens
    const actions = ['create', 'update', 'delete'];
    const action = actions[Math.floor(Math.random() * actions.length)];
    
    // Note: These will fail without auth, but still test the system
    if (action === 'create') {
      await makeRequest('POST', '/blogs/create', {
        title: `Load Test Blog ${userId}`,
        content: 'This is a load test blog post',
        tags: ['test']
      });
    }
  }
}

// Run concurrent users
async function runLoadTest() {
  console.log('🚀 Starting Load Test...\n');
  console.log('Configuration:');
  console.log(`  - Total Users: ${CONFIG.totalUsers.toLocaleString()}`);
  console.log(`  - Concurrent Users: ${CONFIG.concurrentUsers}`);
  console.log(`  - Test Duration: ${CONFIG.testDurationMinutes} minutes\n`);
  
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
    
    // Small delay between batches to not overwhelm the system instantly
    await new Promise(resolve => setTimeout(resolve, 100));
  }
  
  clearInterval(statusInterval);
  metrics.endTime = Date.now();
  
  printFinalReport();
}

function printStatus() {
  const elapsed = ((Date.now() - metrics.startTime) / 1000).toFixed(1);
  const rps = (metrics.totalRequests / (elapsed || 1)).toFixed(2);
  const successRate = ((metrics.successfulRequests / metrics.totalRequests) * 100).toFixed(2);
  
  console.log(`⏱️  Elapsed: ${elapsed}s | Requests: ${metrics.totalRequests} | RPS: ${rps} | Success: ${successRate}%`);
}

function calculatePercentile(arr, percentile) {
  const sorted = [...arr].sort((a, b) => a - b);
  const index = Math.ceil((percentile / 100) * sorted.length) - 1;
  return sorted[index] || 0;
}

function printFinalReport() {
  console.log('\n' + '='.repeat(70));
  console.log('LOAD TEST RESULTS');
  console.log('='.repeat(70));
  
  const duration = (metrics.endTime - metrics.startTime) / 1000;
  const rps = metrics.totalRequests / duration;
  const successRate = (metrics.successfulRequests / metrics.totalRequests) * 100;
  
  console.log('\n📊 Overall Statistics:');
  console.log(`  Total Requests: ${metrics.totalRequests.toLocaleString()}`);
  console.log(`  Successful: ${metrics.successfulRequests.toLocaleString()} (${successRate.toFixed(2)}%)`);
  console.log(`  Failed: ${metrics.failedRequests.toLocaleString()}`);
  console.log(`  Duration: ${duration.toFixed(2)}s`);
  console.log(`  Requests/sec: ${rps.toFixed(2)}`);
  
  console.log('\n⏱️  Response Times:');
  const avg = metrics.responseTimes.reduce((a, b) => a + b, 0) / metrics.responseTimes.length;
  console.log(`  Average: ${avg.toFixed(2)}ms`);
  console.log(`  Min: ${Math.min(...metrics.responseTimes).toFixed(2)}ms`);
  console.log(`  Max: ${Math.max(...metrics.responseTimes).toFixed(2)}ms`);
  console.log(`  P50: ${calculatePercentile(metrics.responseTimes, 50).toFixed(2)}ms`);
  console.log(`  P95: ${calculatePercentile(metrics.responseTimes, 95).toFixed(2)}ms`);
  console.log(`  P99: ${calculatePercentile(metrics.responseTimes, 99).toFixed(2)}ms`);
  
  console.log('\n🎯 Per-Endpoint Performance:');
  for (const [endpoint, stats] of Object.entries(metrics.endpoints)) {
    const avgTime = stats.totalTime / stats.count;
    const errorRate = (stats.errors / stats.count) * 100;
    
    console.log(`\n  ${endpoint}:`);
    console.log(`    Requests: ${stats.count}`);
    console.log(`    Avg Time: ${avgTime.toFixed(2)}ms`);
    console.log(`    Min/Max: ${stats.minTime.toFixed(2)}ms / ${stats.maxTime.toFixed(2)}ms`);
    console.log(`    Errors: ${stats.errors} (${errorRate.toFixed(2)}%)`);
  }
  
  if (Object.keys(metrics.errors).length > 0) {
    console.log('\n❌ Errors:');
    for (const [error, count] of Object.entries(metrics.errors)) {
      console.log(`  ${error}: ${count}`);
    }
  }
  
  console.log('\n' + '='.repeat(70));
  
  // Performance assessment
  console.log('\n📈 Performance Assessment:');
  if (avg < 50) {
    console.log('  ✅ EXCELLENT - Average response time < 50ms');
  } else if (avg < 100) {
    console.log('  ✅ GOOD - Average response time < 100ms');
  } else if (avg < 200) {
    console.log('  ⚠️  ACCEPTABLE - Average response time < 200ms');
  } else {
    console.log('  ❌ POOR - Average response time > 200ms - Optimization needed!');
  }
  
  if (successRate > 99) {
    console.log('  ✅ EXCELLENT - Success rate > 99%');
  } else if (successRate > 95) {
    console.log('  ✅ GOOD - Success rate > 95%');
  } else {
    console.log('  ❌ POOR - Too many failed requests!');
  }
  
  console.log('\n' + '='.repeat(70));
}

// Run the test
runLoadTest().catch(console.error);