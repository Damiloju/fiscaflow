# FiscaFlow Development Prompts

This document tracks all prompts and instructions exchanged during the development of FiscaFlow, a personal finance tracker backend built with Go and OpenTelemetry.

## 📋 **Conversation Timeline**

### **Initial Setup & Cursor Rules**

#### **Prompt 1: Cursor Rules Organization**
**Full User Prompt**: "Look at and suggest changes to Cursor rules in the .cursor/rules directory"

**Context**: User wanted to review and improve existing Cursor rules
**Outcome**: 
- Found 2 rule files: `conventional-commits.mdc` and `cursor-rules-location.mdc`
- Suggested format improvements
- Created 3 new rules: `code-quality.mdc`, `security-best-practices.mdc`
- Added 3 more rules: `testing-standards.mdc`, `documentation-standards.mdc`, `performance-guidelines.mdc`

#### **Prompt 2: Rule Organization**
**Full User Prompt**: "Organize the rules into folders"

**Context**: User wanted better organization of the Cursor rules
**Outcome**:
- Created 4 folders: `workflow/`, `development/`, `security/`, `quality/`
- Moved rules to appropriate folders
- Created README files for each folder and main rules directory

### **Product Development Phase**

#### **Prompt 3: High-Value Product Ideas**
**Full User Prompt**: "I want to vibe code a high-value product backend. Give me ideas"

**Context**: User wanted product ideas for backend development
**Outcome**: Provided 12 product ideas categorized by complexity and market potential:
- **Top Recommendations**: Personal Finance Tracker, Freelancer Time Tracker, AI-powered Writing Assistant
- **Medium Complexity**: E-commerce Analytics, Real Estate Management, Healthcare Appointment System
- **High Complexity**: Supply Chain Optimization, IoT Fleet Management, Blockchain-based Voting System

#### **Prompt 4: Product Selection**
**Full User Prompt**: "We are going to create an Personal Finance Tracker using

Find a fancy name for it.

It uses the opentelemetry  for logging, metrics and telemetry.
All operations have appropriate tracing on them that can be used to troubleshoot the application.


IMPORTANT: Write up the specifications into the "specs/" folder with each domain topic (including technical topic) as a seperate markdown file. Create a "SPECS.md" in the root of the directory which is an overview document that contains a table that links to all the specs."

**Context**: User selected Personal Finance Tracker as the product to build
**Outcome**:
- Created product name "FiscaFlow"
- Generated comprehensive specifications in `specs/` folder
- Created `SPECS.md` overview and detailed specs for:
  - OpenTelemetry integration
  - System architecture
  - Database design

#### **Prompt 5: Technology Clarification**
**Full User Prompt**: "The backend should be built with Go, not Node.js"

**Context**: User clarified the technology stack preference
**Outcome**: Updated all specifications to use Go:
- Go-specific frameworks (Gin, GORM)
- Go libraries (Zap for logging)
- Go OpenTelemetry SDK
- Go code examples for service setup, tracing, metrics, logging
- Updated deployment configurations

#### **Prompt 6: Architecture Decision**
**Full User Prompt**: "Use a well structured monolith for now. Allowing for easy porting to microservice architecture when needed"

**Context**: User wanted monolith-first approach with migration path
**Outcome**: Updated system architecture to:
- Well-structured monolith with clean domain separation
- Domain-Driven Design principles
- Event-driven communication within monolith
- 4-phase migration path to microservices
- Go-specific project structure
- Single deployment with Kubernetes

#### **Prompt 7: Prompt Documentation**
**Full User Prompt**: "Create a Prompt.md file in the root directory and organize all the prompts in the chat into it. Update it as we build along with every prompt given. Create a cursor rule to make sure instructions are followed"

**Context**: User wanted documentation of all prompts and a rule to ensure instruction compliance
**Outcome**: 
- Created this Prompt.md file
- Created Cursor rule for instruction compliance

#### **Prompt 8: Full Prompts Request**
**Full User Prompt**: "I want the full prompts"

**Context**: User wanted the complete, exact prompts as they were given, not just summaries
**Outcome**: 
- Updated Prompt.md to include full, exact prompts
- Maintained chronological order and context

#### **Prompt 9: Project Initialization and Structure**
**Full User Prompt**: "- Initialize the Go project and create the directory structure
- Set up the core dependencies and go.mod file
- Follow @system-architecture.md 
- Don't forget to update prompt"

**Context**: User requested to start actual implementation by initializing the Go project, setting up the directory structure and dependencies as specified in system-architecture.md, and to update the prompt log accordingly.
**Outcome**:
- Initialized Go module (`go mod init fiscaflow`)
- Created directory structure as per system-architecture.md
- Installed all core dependencies (Gin, GORM, Zap, OpenTelemetry, etc.)
- Created initial main.go and config.go files
- Set up logging system with Zap
- Updated Prompt.md with this instruction

#### **Prompt 10: Git Ignore and Basic Readme**
**Full User Prompt**: "Create git ignore and Basic Readme"

**Context**: User requested to add a .gitignore file for the Go project and a basic README.md describing the project, stack, and setup.
**Outcome**:
- Created .gitignore with Go, Docker, editor, and OS-specific ignores
- Created README.md with project description, features, tech stack, setup, and references to specs and Prompt.md
- Noted the need to follow the conventional commit rule for these changes

## 🎯 **Key Decisions Made**

### **Technology Stack**
- **Backend**: Go (Gin, GORM, Zap)
- **Database**: PostgreSQL (primary), Redis (cache), Elasticsearch (search)
- **File Storage**: MinIO
- **Message Queue**: RabbitMQ
- **Observability**: OpenTelemetry
- **Deployment**: Docker + Kubernetes

### **Architecture Approach**
- **Pattern**: Well-structured monolith with domain separation
- **Design**: Domain-Driven Design with Clean Architecture
- **Communication**: Event-driven within monolith
- **Migration**: Clear path to microservices when needed

### **Product Features**
- **Core**: Transaction management, categorization, analytics
- **Advanced**: ML-powered categorization, budget tracking, goal setting
- **Security**: JWT authentication, RBAC, encryption
- **Observability**: Full OpenTelemetry integration

## 📝 **Development Instructions**

### **Current State**
- ✅ Product concept defined (FiscaFlow - Personal Finance Tracker)
- ✅ Technology stack selected (Go-based monolith)
- ✅ Architecture designed (well-structured monolith with migration path)
- ✅ Specifications created (OpenTelemetry, System Architecture, Database Design)
- ✅ Prompt documentation system created
- 🔄 **Next**: Implementation phase

### **Pending Tasks**
1. **Project Setup**: Initialize Go modules and project structure
2. **Core Domains**: Implement User and Transaction domains
3. **API Layer**: Create RESTful APIs with Gin
4. **Database**: Set up PostgreSQL with migrations
5. **Observability**: Integrate OpenTelemetry
6. **Testing**: Unit and integration tests
7. **Deployment**: Docker and Kubernetes setup

## 🔄 **Instruction Compliance Rule**

A Cursor rule has been created (`.cursor/rules/instruction-compliance.mdc`) to ensure:
- All prompts are documented in this file with full, exact text
- Instructions are followed precisely
- Progress is tracked systematically
- Decisions are recorded with rationale

## 📋 **Full Conversation Context**

### **Initial Discovery**
The conversation began with the user asking about a mysterious `.history` file created by VS Code, which was identified as being from the Local History extension. This led to a discussion about Cursor rules and their organization.

### **Development Journey**
The conversation evolved from Cursor rules organization to product ideation, technology selection, architecture decisions, and finally to creating a comprehensive documentation system for tracking all prompts and instructions.

### **Documentation Standards**
- All prompts are recorded with their full, exact text
- Context and outcomes are documented for each prompt
- Technical decisions and rationale are preserved
- Progress tracking is maintained systematically

---

*This document will be updated with each new prompt and instruction as we continue building FiscaFlow. All prompts are recorded with their full, exact text as given by the user.* 