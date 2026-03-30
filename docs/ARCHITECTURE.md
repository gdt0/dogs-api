# Dogs API - Technical Documentation

**Version:** 1.0.0  
**Last Updated:** 2026-03-30

---

## 📐 System Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              CLIENT LAYER                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │   Browser   │  │   Mobile    │  │    API      │  │   cURL/     │    │
│  │  (Web UI)   │  │    App      │  │   Client    │  │   Postman   │    │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘    │
└─────────┼────────────────┼────────────────┼────────────────┼────────────┘
          │                │                │                │
          │    HTTP/REST   │                │                │
          ▼                ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           TRANSPORT LAYER                                │
│                        localhost:8080 / HTTPS                            │
│                                                                            │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                        Gin HTTP Server                           │   │
│   │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │   │
│   │  │   CORS   │  │  Logger  │  │ Recovery │  │ Static Files (/, │ │   │
│   │  │ Middleware│ │Middleware│  │Middleware│  │ /static/*, /api/*│ │   │
│   │  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘ │   │
│   └─────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          APPLICATION LAYER                              │
│                                                                            │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                      API Routes (/api/*)                        │   │
│   │                                                                   │   │
│   │   GET /api/dogs           → listAllDogs()                        │   │
│   │   GET /api/dogs/:breed    → getBreed()                          │   │
│   │   POST /api/dogs          → createBreed()                       │   │
│   │   PUT /api/dogs/:breed    → updateBreed()                       │   │
│   │   DELETE /api/dogs/:breed → deleteBreed()                       │   │
│   │                                                                   │   │
│   └─────────────────────────────────────────────────────────────────┘   │
│                                    │                                     │
│                                    ▼                                     │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                    Service Functions                             │   │
│   │                                                                   │   │
│   │   loadDogs()   → Read JSON from disk                             │   │
│   │   saveDogs()   → Write JSON to disk                              │   │
│   │                                                                   │   │
│   └─────────────────────────────────────────────────────────────────┘   │
│                                    │                                     │
└────────────────────────────────────┼────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                            DATA LAYER                                    │
│                                                                            │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                      dogs.json                                   │   │
│   │                                                                   │   │
│   │  {                                                               │   │
│   │    "affenpinscher": [],                                          │   │
│   │    "bulldog": ["boston", "french"],                              │   │
│   │    "labrador": ["golden", "black"]                              │   │
│   │  }                                                               │   │
│   │                                                                   │   │
│   └─────────────────────────────────────────────────────────────────┘   │
│                                                                            │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 🔄 Request/Response Flow

### Create Breed Flow

```
┌────────┐         ┌─────────┐         ┌─────────┐         ┌─────────┐
│ Client │         │   Gin   │         │ Service │         │  File   │
└───┬────┘         └────┬────┘         └────┬────┘         └────┬────┘
    │                   │                   │                   │
    │ POST /api/dogs    │                   │                   │
    │ {breed, varieties}                   │                   │
    │ ───────────────► │                   │                   │
    │                   │                   │                   │
    │                   │ Validate input    │                   │
    │                   │ ─────────────────►│                   │
    │                   │                   │                   │
    │                   │                   │ loadDogs()        │
    │                   │                   │ ───────────────► │
    │                   │                   │   Read file      │
    │                   │                   │ ◄─────────────── │
    │                   │                   │   Return data    │
    │                   │                   │                   │
    │                   │                   │ Check if exists? │
    │                   │                   │ ───────────────► │
    │                   │                   │ ◄─────────────── │
    │                   │                   │                   │
    │                   │                   │ Append new breed  │
    │                   │                   │ to in-memory map  │
    │                   │                   │                   │
    │                   │                   │ saveDogs()        │
    │                   │                   │ ───────────────► │
    │                   │                   │   Write file     │
    │                   │                   │ ◄─────────────── │
    │                   │                   │                   │
    │  201 Created      │                   │                   │
    │ ◄─────────────── │                   │                   │
    │                   │                   │                   │
```

### Delete Breed Flow

```
┌────────┐         ┌─────────┐         ┌─────────┐         ┌─────────┐
│ Client │         │   Gin   │         │ Service │         │  File   │
└───┬────┘         └────┬────┘         └────┬────┘         └────┬────┘
    │                   │                   │                   │
    │ DELETE /api/dogs/ │                   │                   │
    │    labrador        │                   │                   │
    │ ───────────────► │                   │                   │
    │                   │                   │                   │
    │                   │ loadDogs()        │                   │
    │                   │ ───────────────► │                   │
    │                   │                   │ Read file        │
    │                   │ ◄─────────────── │                   │
    │                   │   Return data     │                   │
    │                   │                   │                   │
    │                   │ Breed exists?     │                   │
    │                   │ ───────────────► │                   │
    │                   │                   │                   │
    │                   │   delete from map │                   │
    │                   │                   │                   │
    │                   │ saveDogs()        │                   │
    │                   │ ───────────────► │                   │
    │                   │                   │ Write file       │
    │                   │ ◄─────────────── │                   │
    │                   │                   │                   │
    │  200 OK           │                   │                   │
    │ ◄─────────────── │                   │                   │
    │                   │                   │                   │
```

---

## 📊 Data Model

### Dogs Entity

```
┌─────────────────────────────────────────────────────────────┐
│                         dogs.json                           │
├─────────────────────────────────────────────────────────────┤
│  Key (string)          │  Value (array of strings)         │
├─────────────────────────┼───────────────────────────────────┤
│  "affenpinscher"       │  []                               │
│  "bulldog"             │  ["boston", "french"]             │
│  "labrador"            │  ["golden", "black", "chocolate"] │
│  "poodle"              │  ["miniature", "standard", "toy"] │
└─────────────────────────┴───────────────────────────────────┘
```

### In-Memory Representation (Go)

```go
type Dogs map[string][]string

// Example:
dogs := Dogs{
    "affenpinscher": {},
    "bulldog":       {"boston", "french"},
    "labrador":      {"golden", "black", "chocolate"},
}
```

---

## 🎯 API Endpoints Detail

### GET /api/dogs

**Description:** Retrieve all dog breeds

**Request:**
```
GET /api/dogs
```

**Response (200 OK):**
```json
{
  "affenpinscher": [],
  "african": [],
  "airedale": [],
  "akita": [],
  "basenji": [],
  "beagle": [],
  "bulldog": ["boston", "french"],
  "labrador": ["golden", "black"],
  ...
}
```

**Flow:**
```
1. Receive request
2. Call loadDogs()
3. Return JSON response
```

---

### GET /api/dogs/:breed

**Description:** Retrieve a specific breed

**Request:**
```
GET /api/dogs/labrador
```

**Response (200 OK):**
```json
{
  "breed": "labrador",
  "varieties": ["golden", "black"]
}
```

**Response (404 Not Found):**
```json
{
  "error": "Breed not found"
}
```

**Flow:**
```
1. Extract breed parameter (lowercase)
2. Call loadDogs()
3. Check if breed exists in map
4. Return 200 with data OR 404 with error
```

---

### POST /api/dogs

**Description:** Create a new breed

**Request:**
```
POST /api/dogs
Content-Type: application/json

{
  "breed": "newbreed",
  "varieties": ["variety1", "variety2"]
}
```

**Response (201 Created):**
```json
{
  "breed": "newbreed",
  "varieties": ["variety1", "variety2"]
}
```

**Response (409 Conflict):**
```json
{
  "error": "Breed already exists"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Key: 'Input.breed' Error:Field validation failed"
}
```

**Flow:**
```
1. Parse JSON body
2. Validate breed field required
3. Normalize breed to lowercase
4. Call loadDogs()
5. Check if breed already exists
6. Add to map
7. Call saveDogs()
8. Return 201 with new breed
```

---

### PUT /api/dogs/:breed

**Description:** Update a breed's varieties

**Request:**
```
PUT /api/dogs/labrador
Content-Type: application/json

{
  "varieties": ["yellow", "chocolate"]
}
```

**Response (200 OK):**
```json
{
  "breed": "labrador",
  "varieties": ["yellow", "chocolate"]
}
```

**Response (404 Not Found):**
```json
{
  "error": "Breed not found"
}
```

**Flow:**
```
1. Extract breed parameter
2. Parse JSON body
3. Call loadDogs()
4. Check if breed exists
5. Update varieties in map
6. Call saveDogs()
7. Return 200 with updated data
```

---

### DELETE /api/dogs/:breed

**Description:** Delete a breed

**Request:**
```
DELETE /api/dogs/labrador
```

**Response (200 OK):**
```json
{
  "message": "Breed 'labrador' deleted"
}
```

**Response (404 Not Found):**
```json
{
  "error": "Breed not found"
}
```

**Flow:**
```
1. Extract breed parameter
2. Call loadDogs()
3. Check if breed exists
4. Delete from map
5. Call saveDogs()
6. Return 200 with confirmation
```

---

## 🏗️ Component Design

### HTTP Server (main.go)

```
┌─────────────────────────────────────────────────────────────────┐
│                         Gin Engine                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Middleware Stack (top to bottom):                              │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ Logger        → Logs all requests with status/time      │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │ Recovery      → Catches panics, returns 500            │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │ CORS          → Adds CORS headers, handles preflight    │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                  │
│  Routes:                                                        │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ GET  /              → Serve index.html                   │   │
│  │ GET  /static/*      → Serve static files                 │   │
│  │ GET  /api/dogs      → listAllDogs()                      │   │
│  │ GET  /api/dogs/:breed → getBreed()                       │   │
│  │ POST /api/dogs      → createBreed()                      │   │
│  │ PUT  /api/dogs/:breed → updateBreed()                    │   │
│  │ DELETE /api/dogs/:breed→ deleteBreed()                   │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Service Layer

```
┌─────────────────────────────────────────────────────────────────┐
│                       Service Functions                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  loadDogs() ─────────────────────────────────────────────────   │
│  │                                                               │
│  ├─ Read file "dogs.json"                                       │
│  ├─ Unmarshal JSON to Dogs map                                  │
│  ├─ Return Dogs or error                                        │
│  └─ Called by all read operations                              │
│                                                                  │
│  saveDogs(dogs Dogs) ────────────────────────────────────────  │
│  │                                                               │
│  ├─ Marshal Dogs map to JSON (pretty printed)                   │
│  ├─ Write to "dogs.json"                                        │
│  ├─ Set permissions to 0644                                    │
│  └─ Return error or nil                                        │
│      Called by all write operations                             │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📈 State Diagram

```
                    ┌─────────────────┐
                    │   Application   │
                    │     Starts      │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │   Load dogs.json │
                    │   into memory   │
                    └────────┬────────┘
                             │
                             ▼
              ┌──────────────────────────────┐
              │        Running State          │
              │  ┌─────────────────────────┐  │
              │  │ Memory: Dogs map        │  │
              │  │ Disk: dogs.json         │  │
              │  │ Serving HTTP requests  │  │
              │  └─────────────────────────┘  │
              └──────────────┬───────────────┘
                             │
            ┌────────────────┼────────────────┐
            │                │                │
            ▼                ▼                ▼
   ┌────────────────┐ ┌────────────┐ ┌──────────────┐
   │   GET Request  │ │  POST Req  │ │ DELETE Req   │
   └───────┬────────┘ └─────┬──────┘ └──────┬───────┘
           │                 │               │
           ▼                 ▼               ▼
   ┌────────────────┐ ┌────────────┐ ┌──────────────┐
   │ Return from    │ │ Validate   │ │ Load dogs    │
   │ memory map     │ │ Check if   │ │ Delete from  │
   │ No state       │ │ exists     │ │ memory map   │
   │ change         │ │            │ │              │
   └────────────────┘ └─────┬──────┘ └──────┬───────┘
                             │                │
                    ┌────────┴────────┐       │
                    ▼                 ▼       ▼
             ┌────────────┐     ┌──────────┐ ┌──────────────┐
             │ Exists?    │     │ Save dogs│ │ Save dogs    │
             │ Yes → 409  │     │ to disk  │ │ to disk     │
             │ No → Add   │     └────┬─────┘ └──────┬───────┘
             └─────┬──────┘          │              │
                   │                 ▼              │
                   ▼          ┌───────────┐         │
              ┌────────┐      │ Return    │         │
              │ Save   │      │ 201       │         │
              │ to disk│      └───────────┘         │
              └───┬────┘                             │
                  │                                  │
                  └──────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   Back to       │
                    │   Running       │
                    └─────────────────┘
```

---

## 🔌 Middleware Chain

```
Incoming Request
       │
       ▼
┌──────────────────┐
│     CORS         │
│  - Set headers   │
│  - Handle OPT    │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│     Logger       │
│  - Log request   │
│  - Log response  │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│    Recovery      │
│  - Catch panic   │
│  - Return 500    │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│    Router        │
│  - Match route   │
│  - Call handler  │
└────────┬─────────┘
         │
         ▼
   Handler Function
   (API endpoint)
```

---

## 🗂️ File Structure

```
dogs-api/
│
├── main.go                 # Application entry point
│   ├── main()             # Server initialization
│   ├── cors()             # CORS middleware factory
│   ├── loadDogs()         # Read JSON file
│   └── saveDogs()         # Write JSON file
│
├── index.html             # Web UI (served at /)
│   ├── HTML structure
│   ├── CSS styles
│   └── Vanilla JS (API client)
│
├── dogs.json              # Data file
│
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
│
├── Dockerfile             # Docker image build
├── docker-entrypoint.sh   # Container startup
│
└── README.md              # Project documentation
```

---

## 🔒 Error Handling

```
┌─────────────────────────────────────────────────────────────────┐
│                     Error Handling Strategy                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  File Errors:                                                    │
│  ├─ dogs.json missing    → Return empty Dogs{} (first run)     │
│  ├─ dogs.json corrupted   → Return 500 with error message       │
│  └─ dogs.json permission  → Return 500 with error message       │
│                                                                  │
│  Validation Errors:                                             │
│  ├─ Missing breed field  → Return 400 with validation error    │
│  └─ Invalid JSON body     → Return 400 with parse error          │
│                                                                  │
│  Business Logic Errors:                                          │
│  ├─ Breed not found      → Return 404 with "Breed not found"     │
│  └─ Breed already exists → Return 409 with "Breed already exists"│
│                                                                  │
│  Server Errors:                                                  │
│  └─ Any unhandled panic  → Recovery middleware catches, returns │
│                             500 with "Internal Server Error"      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🌐 CORS Configuration

```
┌─────────────────────────────────────────────────────────────────┐
│                     CORS Headers                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Access-Control-Allow-Origin: *                                 │
│  → Allows requests from any origin                              │
│                                                                  │
│  Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS  │
│  → Allowed HTTP methods                                         │
│                                                                  │
│  Access-Control-Allow-Headers: Content-Type                     │
│  → Allowed request headers                                      │
│                                                                  │
│  Preflight (OPTIONS) requests:                                  │
│  → Return 204 No Content with headers                           │
│  → Do not run route handlers                                    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📦 Dependencies

```
github.com/gin-gonic/gin v1.9.1
│
├── github.com/gin-contrib/sse v0.1.0          # SSE support
├── github.com/go-playground/validator/v10     # Request validation
├── github.com/mattn/go-isatty v0.0.19        # Terminal detection
├── github.com/ugorji/go/codec v1.2.11        # msgpack/codec
├── github.com/bytedance/sonic v1.9.1         # JSON serialization
├── github.com/goccy/go-json v0.10.2         # JSON parser
├── github.com/json-iterator/go v1.1.12       # JSON iterator
├── github.com/pelletier/go-toml/v2 v2.0.8   # TOML parser
├── golang.org/x/net v0.10.0                 # Networking
├── golang.org/x/sys v0.8.0                  # System calls
├── golang.org/x/text v0.9.0                 # Text processing
└── golang.org/x/crypto v0.9.0              # Cryptography
```

---

## 🚀 Deployment Options

### 1. Direct Binary

```
┌─────────────────────────────────────────────────────────────────┐
│                         VPS/Server                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐       │
│  │   Source     │    │   Build     │    │   Run       │       │
│  │   Code       │───►│   Binary    │───►│   Server    │       │
│  └─────────────┘    └─────────────┘    └──────┬──────┘       │
│                                                  │               │
│                                                  ▼               │
│                                           ┌─────────────┐       │
│                                           │  Port 8080  │       │
│                                           └──────┬──────┘       │
│                                                  │               │
└──────────────────────────────────────────────────┼───────────────┘
                                                   │
                                                   ▼
                                            ┌─────────────┐
                                            │   nginx     │
                                            │  (reverse   │
                                            │   proxy)    │
                                            └─────────────┘
```

### 2. Docker

```
┌─────────────────────────────────────────────────────────────────┐
│                      Docker Container                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Alpine Linux                          │   │
│  │  ┌─────────────────────────────────────────────────┐   │   │
│  │  │              Go Binary (server)                  │   │   │
│  │  │  ┌─────────────────────────────────────────┐    │   │   │
│  │  │  │              Gin HTTP Server              │    │   │   │
│  │  │  │  ┌───────────────────────────────────┐   │    │   │   │
│  │  │  │  │  CORS │ Logger │ Recovery │ Routes│   │    │   │   │
│  │  │  │  └───────────────────────────────────┘   │    │   │   │
│  │  │  └─────────────────────────────────────────┘    │   │   │
│  │  └─────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────┘   │
│                            │                                     │
│                            ▼ (port 8080)                        │
│                    ┌───────────────┐                           │
│                    │  Host Network  │                           │
│                    └───────────────┘                           │
└─────────────────────────────────────────────────────────────────┘
```

### 3. Localtunnel (Development)

```
┌──────────────┐      ┌──────────────┐      ┌──────────────────┐
│   Browser    │─────►│  Localtunnel  │─────►│   Public URL     │
│              │◄─────│   Server     │◄─────│  lee-dogs-api    │
│              │      │              │      │  .loca.lt        │
└──────────────┘      └──────────────┘      └──────────────────┘
```

---

## ⚡ Performance Characteristics

```
┌─────────────────────────────────────────────────────────────────┐
│                     Performance Metrics                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Startup Time:     ~100ms                                       │
│  Request Latency:  <10ms (in-memory operations)                │
│  File I/O:         ~1-5ms per operation                          │
│  Memory Usage:     ~10-50MB (depends on data size)              │
│  Binary Size:      ~11MB (compressed Alpine: ~5MB)              │
│                                                                  │
│  Concurrency:      Gin handles concurrent requests              │
│  Connection:       Keep-alive enabled by default                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🧪 Testing Strategy

```
┌─────────────────────────────────────────────────────────────────┐
│                       API Testing                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  cURL Examples:                                                  │
│                                                                  │
│  # Create                                                         │
│  curl -X POST http://localhost:8080/api/dogs \                  │
│    -H "Content-Type: application/json" \                        │
│    -d '{"breed":"testdog","varieties":["a","b"]}'               │
│                                                                  │
│  # Read                                                           │
│  curl http://localhost:8080/api/dogs                             │
│  curl http://localhost:8080/api/dogs/testdog                    │
│                                                                  │
│  # Update                                                         │
│  curl -X PUT http://localhost:8080/api/dogs/testdog \            │
│    -H "Content-Type: application/json" \                        │
│    -d '{"varieties":["c","d"]}'                                 │
│                                                                  │
│  # Delete                                                        │
│  curl -X DELETE http://localhost:8080/api/dogs/testdog           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

**Document Version:** 1.0  
**Architecture Style:** RESTful JSON API  
**Frontend:** Single Page Application (SPA)  
**Persistence:** File-based (JSON)  
**Containerization:** Docker
