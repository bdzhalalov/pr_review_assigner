# Pull Request reviewer Assigner

The service is designed for managing team pull requests. It allows you to create PRs, reassign reviewers, and merge PRs. It also allows you to manage teams (creating and receiving teams) and user activity.

---
## Getting started

### Prerequisites
- **Docker** with Docker Compose

### Start up project

- Go to directory with cloned project
- Use command `make start`

#### For subsequent launches of the application, it is enough to use the command `make run`

---

## Load Testing
 
**Tool:** k6  
**Scenario:** 3-stage load test with 1–20 virtual users (VUs), total duration ~1 minute.  

### Tested Endpoints

1. `POST /api/teams/add`  
2. `GET /api/teams/get`  
3. `POST /api/users/setIsActive`  
4. `POST /api/pullRequest/create`  
5. `POST /api/pullRequest/merge`  
6. `POST /api/pullRequest/reassign`  
7. `GET /api/users/getReview`

### Results

#### Checks (`check()`)

| Endpoint                     | Total Checks | Passed | Failed |
|-------------------------------|---------------|---------|--------|
| teams/add                     | 790           | 790     | 0      |
| teams/get                     | 790           | 790     | 0      |
| setIsActive                   | 790           | 790     | 0      |
| pr/create                     | 790           | 790     | 0      |
| pr/merge                      | 790           | 790     | 0      |
| pr/reassign                    | 790           | 790     | 0      |
| users/getReview               | 790           | 790     | 0      |

✅ All checks passed successfully. Requests that could return 400/404/409 were marked as expected and are not considered business logic errors.


#### HTTP Metrics

| Metric                       | Value |
|-------------------------------|----------|
| Total Requests                | 5530     |
| Average Response Time (ms)     | 5.19     |
| Min Response Time (ms)        | 1.45     |
| Max Response Time (ms)       | 63.67    |
| 90th Percentile (ms)          | 9.12     |
| 95th Percentile (ms)         | 10.76    |
| HTTP Errors (4xx/5xx)         | 42.82%  |

> Note: The high percentage of HTTP errors is due to test data (duplicate inserts, non-existent resources). This is expected behavior for verifying business logic.

#### Load

- Virtual Users (VUs): 1–20  
- Average Iteration Duration: 1.03 s  
- Maximum Iteration Duration: 1.08 s  

### Conclusion

The system handles up to 20 concurrent users under the current scenarios. The average response time is ~5 ms, demonstrating good performance. All business logic checks passed successfully, and HTTP errors correspond to expected behavior for the test cases.

---

## TODOs

1. Add unit tests
2. Start using transactions for multiple queries
3. Start using Context
4. Add validation to services

