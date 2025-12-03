// Blog API Automated Test Suite
// Run with: node test-api.js
// Make sure your server is running on http://localhost:8080

const BASE_URL = 'http://localhost:8080';

// Test results tracking
const results = {
  passed: [],
  failed: [],
  total: 0
};

// Helper function to make requests
async function makeRequest(method, endpoint, body = null, token = null) {
  const options = {
    method,
    headers: {
      'Content-Type': 'application/json',
    }
  };

  if (token) {
    options.headers['Authorization'] = `Bearer ${token}`;
  }

  if (body) {
    options.body = JSON.stringify(body);
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, options);
  const text = await response.text();
  let data = null;
  
  try {
    data = text ? JSON.parse(text) : null;
  } catch (e) {
    data = text;
  }

  return { status: response.status, data };
}

// Test helper
function test(name, passed, details = '') {
  results.total++;
  if (passed) {
    results.passed.push(name);
    console.log(`✅ ${name}`);
  } else {
    results.failed.push({ name, details });
    console.log(`❌ ${name}`);
    if (details) console.log(`   ${details}`);
  }
}

// Generate random username to avoid conflicts
const randomId = Math.floor(Math.random() * 100000);
const testUser = {
  username: `testuser_${randomId}`,
  password: 'TestPass123!'
};

let authToken = null;
let userId = null;
let blogId = null;

async function runTests() {
  console.log('🚀 Starting API Tests...\n');
  console.log('='.repeat(50));
  console.log('AUTHENTICATION TESTS');
  console.log('='.repeat(50));

  // Test 1: Register (FIXED: Accept 201 status)
  try {
    const res = await makeRequest('POST', '/register', testUser);
    test(
      '1. POST /register - Create new user',
      (res.status === 200 || res.status === 201) && res.data.token && res.data.username,
      (res.status !== 200 && res.status !== 201) ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
    );
    
    if (res.data.token) {
      authToken = res.data.token;
      try {
        const payload = JSON.parse(atob(authToken.split('.')[1]));
        userId = payload.user_id;
      } catch (e) {
        console.log('   Note: Could not decode token for user_id');
      }
    }
  } catch (error) {
    test('1. POST /register - Create new user', false, error.message);
  }

  // Test 2: Login
  try {
    const res = await makeRequest('POST', '/login', testUser);
    test(
      '2. POST /login - Login with credentials',
      res.status === 200 && res.data.token,
      res.status !== 200 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
    );
    
    if (res.data.token && !authToken) {
      authToken = res.data.token;
    }
  } catch (error) {
    test('2. POST /login - Login with credentials', false, error.message);
  }

  // Test 3: Login with wrong password
  try {
    const res = await makeRequest('POST', '/login', {
      username: testUser.username,
      password: 'wrongpassword'
    });
    test(
      '3. POST /login - Reject invalid credentials',
      res.status === 401 || res.status === 400,
      res.status === 200 ? 'Should have rejected wrong password!' : ''
    );
  } catch (error) {
    test('3. POST /login - Reject invalid credentials', false, error.message);
  }

  console.log('\n' + '='.repeat(50));
  console.log('USER ROUTES TESTS');
  console.log('='.repeat(50));

  // Test 4: Get My Profile (protected)
  try {
    console.log(`\nDEBUG: Using token: ${authToken?.substring(0, 50)}...`);
    const res = await makeRequest('GET', '/users/me', null, authToken);
    test(
      '4. GET /users/me - Get authenticated user profile',
      res.status === 200 && res.data.username === testUser.username,
      res.status !== 200 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
    );
    
    if (res.data && res.data.id) {
      userId = res.data.id;
    }
  } catch (error) {
    test('4. GET /users/me - Get authenticated user profile', false, error.message);
  }

  // Test 5: Get My Profile without token
  try {
    const res = await makeRequest('GET', '/users/me', null, null);
    test(
      '5. GET /users/me - Reject unauthenticated request',
      res.status === 401,
      res.status !== 401 ? 'Should require authentication!' : ''
    );
  } catch (error) {
    test('5. GET /users/me - Reject unauthenticated request', false, error.message);
  }

  // Test 6: Get User Profile by ID (NEW: Check for blog_count)
  if (userId) {
    try {
      const res = await makeRequest('GET', `/users/${userId}`);
      test(
        '6. GET /users/{id} - Get user profile by ID',
        res.status === 200 && res.data.username === testUser.username && 'blog_count' in res.data,
        res.status !== 200 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
      );
    } catch (error) {
      test('6. GET /users/{id} - Get user profile by ID', false, error.message);
    }
  } else {
    test('6. GET /users/{id} - Get user profile by ID', false, 'No userId available');
  }

  console.log('\n' + '='.repeat(50));
  console.log('BLOG ROUTES TESTS');
  console.log('='.repeat(50));

  // Test 7: Create Blog (protected)
  try {
    const blogData = {
      title: 'Test Blog Post',
      content: 'This is a test blog post created by the automated test suite.',
      tags: ['test', 'automation']
    };
    const res = await makeRequest('POST', '/blogs/create', blogData, authToken);
    test(
      '7. POST /blogs/create - Create new blog',
      res.status === 200 || res.status === 201,
      res.status !== 200 && res.status !== 201 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
    );
    
    if (res.data && res.data.id) {
      blogId = res.data.id;
    }
  } catch (error) {
    test('7. POST /blogs/create - Create new blog', false, error.message);
  }

  // Test 8: Create Blog without token
  try {
    const res = await makeRequest('POST', '/blogs/create', {
      title: 'Unauthorized Blog',
      content: 'This should fail'
    }, null);
    test(
      '8. POST /blogs/create - Reject unauthenticated request',
      res.status === 401,
      res.status !== 401 ? 'Should require authentication!' : ''
    );
  } catch (error) {
    test('8. POST /blogs/create - Reject unauthenticated request', false, error.message);
  }

  // Test 9: Get My Blogs (protected)
  try {
    const res = await makeRequest('GET', '/blogs/me', null, authToken);
    test(
      '9. GET /blogs/me - Get authenticated user\'s blogs',
      res.status === 200 && Array.isArray(res.data.blogs || res.data),
      res.status !== 200 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
    );
  } catch (error) {
    test('9. GET /blogs/me - Get authenticated user\'s blogs', false, error.message);
  }

  // Test 10: Get Blog by ID (public)
  if (blogId) {
    try {
      const res = await makeRequest('GET', `/blogs/${blogId}`);
      test(
        '10. GET /blogs/{id} - Get blog by ID',
        res.status === 200 && res.data.title === 'Test Blog Post',
        res.status !== 200 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
      );
    } catch (error) {
      test('10. GET /blogs/{id} - Get blog by ID', false, error.message);
    }
  } else {
    test('10. GET /blogs/{id} - Get blog by ID', false, 'No blogId available');
  }

  // Test 11: List All Blogs (public)
  try {
    const res = await makeRequest('GET', '/blogs?page=1&limit=10');
    test(
      '11. GET /blogs - List all blogs with pagination',
      res.status === 200 && (Array.isArray(res.data.blogs) || Array.isArray(res.data)),
      res.status !== 200 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
    );
  } catch (error) {
    test('11. GET /blogs - List all blogs with pagination', false, error.message);
  }

  // Test 12: Search Blogs (public)
  try {
    const res = await makeRequest('GET', '/blogs/search?q=test');
    test(
      '12. GET /blogs/search - Search blogs',
      res.status === 200,
      res.status !== 200 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
    );
  } catch (error) {
    test('12. GET /blogs/search - Search blogs', false, error.message);
  }

  // Test 13: Get User's Blogs (public)
  if (userId) {
    try {
      const res = await makeRequest('GET', `/users/${userId}/blogs`);
      test(
        '13. GET /users/{id}/blogs - Get user\'s blogs',
        res.status === 200 && (Array.isArray(res.data.blogs) || Array.isArray(res.data)),
        res.status !== 200 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
      );
    } catch (error) {
      test('13. GET /users/{id}/blogs - Get user\'s blogs', false, error.message);
    }
  } else {
    test('13. GET /users/{id}/blogs - Get user\'s blogs', false, 'No userId available');
  }

  // NEW Test 14: Update Blog (protected)
  if (blogId) {
    try {
      const updateData = {
        title: 'Updated Test Blog Post',
        content: 'This blog post has been updated by the test suite.',
      };
      const res = await makeRequest('PUT', `/blogs/${blogId}`, updateData, authToken);
      test(
        '14. PUT /blogs/{id} - Update blog',
        res.status === 200 && res.data.title === 'Updated Test Blog Post',
        res.status !== 200 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
      );
    } catch (error) {
      test('14. PUT /blogs/{id} - Update blog', false, error.message);
    }
  } else {
    test('14. PUT /blogs/{id} - Update blog', false, 'No blogId available');
  }

  // NEW Test 15: Update Blog without token
  if (blogId) {
    try {
      const res = await makeRequest('PUT', `/blogs/${blogId}`, {
        title: 'Unauthorized Update'
      }, null);
      test(
        '15. PUT /blogs/{id} - Reject unauthenticated update',
        res.status === 401,
        res.status !== 401 ? 'Should require authentication!' : ''
      );
    } catch (error) {
      test('15. PUT /blogs/{id} - Reject unauthenticated update', false, error.message);
    }
  }

  // Test 16: Delete Blog (protected)
  if (blogId) {
    try {
      const res = await makeRequest('DELETE', `/blogs/${blogId}`, null, authToken);
      test(
        '16. DELETE /blogs/{id} - Delete blog',
        res.status === 200 || res.status === 204,
        res.status !== 200 && res.status !== 204 ? `Status: ${res.status}, Response: ${JSON.stringify(res.data)}` : ''
      );
    } catch (error) {
      test('16. DELETE /blogs/{id} - Delete blog', false, error.message);
    }
  } else {
    test('16. DELETE /blogs/{id} - Delete blog', false, 'No blogId available');
  }

  // Test 17: Delete Blog without token
  if (blogId) {
    try {
      const res = await makeRequest('DELETE', `/blogs/${blogId}`, null, null);
      test(
        '17. DELETE /blogs/{id} - Reject unauthenticated delete',
        res.status === 401,
        res.status !== 401 ? 'Should require authentication!' : ''
      );
    } catch (error) {
      test('17. DELETE /blogs/{id} - Reject unauthenticated delete', false, error.message);
    }
  }

  // Print Summary
  console.log('\n' + '='.repeat(50));
  console.log('TEST SUMMARY');
  console.log('='.repeat(50));
  console.log(`Total Tests: ${results.total}`);
  console.log(`✅ Passed: ${results.passed.length}`);
  console.log(`❌ Failed: ${results.failed.length}`);
  
  if (results.failed.length > 0) {
    console.log('\nFailed Tests:');
    results.failed.forEach(({ name, details }) => {
      console.log(`  - ${name}`);
      if (details) console.log(`    ${details}`);
    });
  }

  console.log('\n' + '='.repeat(50));
  const passRate = ((results.passed.length / results.total) * 100).toFixed(1);
  console.log(`Pass Rate: ${passRate}%`);
  console.log('='.repeat(50));
}

// Run the tests
runTests().catch(console.error);