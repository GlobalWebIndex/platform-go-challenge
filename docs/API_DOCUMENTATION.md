# API Documentation (AI Generated)

## REST API Endpoints

### Audiences

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/audience` | Create a new audience |
| GET | `/audiences` | Get all audiences |
| GET | `/audience/:id` | Get audience by ID |
| PUT | `/audience/:id` | Update audience by ID |
| DELETE | `/audience/:id` | Delete audience by ID |

**Audience Model:**
```json
{
  "id": 1,
  "gender": "Male",
  "birthcountry": "USA",
  "agegroup": "25-34",
  "dailyhours": 5,
  "noofpurchases": 10
}
```

### Charts

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/chart` | Create a new chart |
| GET | `/charts` | Get all charts |
| GET | `/chart/:id` | Get chart by ID |
| PUT | `/chart/:id` | Update chart by ID |
| DELETE | `/chart/:id` | Delete chart by ID |

**Chart Model:**
```json
{
  "id": 1,
  "title": "Sales Chart",
  "xaxistitle": "Months",
  "yaxistitle": "Revenue"
}
```

### Insights

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/insight` | Create a new insight |
| GET | `/insights` | Get all insights |
| GET | `/insight/:id` | Get insight by ID |
| PUT | `/insight/:id` | Update insight by ID |
| DELETE | `/insight/:id` | Delete insight by ID |

**Insight Model:**
```json
{
  "id": 1,
  "text": "Sales increased by 20% in Q4"
}
```

### User Stars

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/userstar` | Create a new user star |
| GET | `/userstars` | Get all user stars |
| GET | `/userstar/:id` | Get user star by ID |
| PUT | `/userstar/:id` | Update user star by ID |
| DELETE | `/userstar/:id` | Delete user star by ID |

**UserStar Model:**
```json
{
  "id": 1,
  "userid": 123,
  "type": "Chart",
  "assetid": 456
}
```

**Type field** must be one of: `"Audience"`, `"Chart"`, or `"Insight"` (capitalized)

---

## GraphQL API

**Endpoint:** `POST /graphql`
**Playground:** `GET /graphql`

### Queries

#### Audiences
```graphql
# Get all audiences
query {
  audiences {
    id
    gender
    birthcountry
    agegroup
    dailyhours
    noofpurchases
  }
}

# Get audience by ID
query {
  audience(id: "1") {
    id
    gender
    birthcountry
  }
}
```

#### Charts
```graphql
# Get all charts
query {
  charts {
    id
    title
    xaxistitle
    yaxistitle
  }
}

# Get chart by ID
query {
  chart(id: "1") {
    id
    title
  }
}
```

#### Insights
```graphql
# Get all insights
query {
  insights {
    id
    text
  }
}

# Get insight by ID
query {
  insight(id: "1") {
    id
    text
  }
}
```

#### User Stars
```graphql
# Get all user stars
query {
  userstars {
    id
    userid
    type
    assetid
  }
}

# Get user star by ID
query {
  userstar(id: "1") {
    id
    userid
    type
    assetid
  }
}

# Get all starred items for a user (aggregated by type)
query {
  userstared(userID: "123") {
    userid
    audience {
      id
      gender
      birthcountry
      agegroup
      dailyhours
      noofpurchases
    }
    chart {
      id
      title
      xaxistitle
      yaxistitle
    }
    insight {
      id
      text
    }
  }
}
```

**Note:** The `userstared` query fetches all user stars for a specific user and returns the full details of each starred asset, grouped by type (audiences, charts, insights).

### Mutations

#### Audiences
```graphql
# Create audience
mutation {
  createAudience(input: {
    gender: "Male"
    birthcountry: "USA"
    agegroup: "25-34"
    dailyhours: 5
    noofpurchases: 10
  }) {
    id
    gender
  }
}

# Update audience
mutation {
  updateAudience(id: "1", input: {
    dailyhours: 6
  }) {
    id
    dailyhours
  }
}

# Delete audience
mutation {
  deleteAudience(id: "1")
}
```

#### Charts
```graphql
# Create chart
mutation {
  createChart(input: {
    title: "Sales Chart"
    xaxistitle: "Months"
    yaxistitle: "Revenue"
  }) {
    id
    title
  }
}

# Update chart
mutation {
  updateChart(id: "1", input: {
    title: "Updated Sales Chart"
  }) {
    id
    title
  }
}

# Delete chart
mutation {
  deleteChart(id: "1")
}
```

#### Insights
```graphql
# Create insight
mutation {
  createInsight(input: {
    text: "Sales increased by 20% in Q4"
  }) {
    id
    text
  }
}

# Update insight
mutation {
  updateInsight(id: "1", input: {
    text: "Sales increased by 25% in Q4"
  }) {
    id
    text
  }
}

# Delete insight
mutation {
  deleteInsight(id: "1")
}
```

#### User Stars
```graphql
# Create user star
mutation {
  createUserStar(input: {
    userid: 123
    type: "Chart"
    assetid: 456
  }) {
    id
    userid
    type
    assetid
  }
}

# Update user star
mutation {
  updateUserStar(id: "1", input: {
    type: "Insight"
    assetid: 789
  }) {
    id
    type
    assetid
  }
}

# Delete user star
mutation {
  deleteUserStar(id: "1")
}
```
