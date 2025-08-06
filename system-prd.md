# Product Requirements Document: LLM Prompt Engineering Instruction Management
System

## Introduction {:#introduction}

This document outlines the requirements for a web application designed to manage
and organize LLM prompt engineering instructions, specifically those intended
for coding agents. The system will automate the ingestion and processing of
long-form markdown-based instructions, break them down into discrete snippets,
and store them in a Firestore database with relevant labels. The web application
will provide users with a search interface to discover and utilize these
snippets, incorporating a crowdsourced rating mechanism for quality assessment
and sorting.

## 1. System Overview {:#1.-system}

The LLM Prompt Engineering Instruction Management System will consist of two
primary components: an automated backend processing pipeline and a user-facing
web application.

## 2. Automated Backend Processing {:#2.-automated}

The backend will be responsible for ingesting, processing, and storing the LLM
prompt engineering instructions.

### 2.1 Ingestion {:#2.1-ingestion}

- **Input Format**: Long-form markdown-based LLM prompt engineering
  instructions.
  - URLs
  - uploaded content
- **Triggers**: Regular ingestion schedule (e.g., daily, weekly) for URLs
  - on upload or edit(for files)

### 2.2 Processing {:#2.2-processing}

- **LLM Integration**: The system will utilize an LLM to break down ingested
  long-form instructions.
- **Snippet Generation**: The LLM will identify and extract discrete, standalone
  instruction snippets from the long-form documents.
- **Labeling**: For each snippet, the LLM will generate a set of relevant
  labels. Examples of labels include:
  - Programming language (e.g., Python, Java, JavaScript)
  - Framework (e.g., React, Angular, Django, Spring Boot)
  - Process (e.g., Git workflow, Agile methodologies)
  - Source code management (e.g., Git, GitHub, GitLab)
  - Other relevant categories ( to define additional categories)
- For each snippet, a gemini embedding model running on vertex will generate an
  embedding
- For each snippet a safety filter is run on content using the appropriate
  vertex API

### 2.3 Storage {:#2.3-storage}

- **Database**: Firestore document database.
- **Data Structure for snippets**: Each snippet will be stored as a document in
  Firestore, containing:
  - Snippet text
  - Associated labels
  - Timestamp of creation
  - Reference to the original long-form instruction (file document store, or
    URL)
  - The embedding value from gemini embedding on the snippet text
  - two counters. One for "thumbs up" one for "thumbs down"
- Storage for sources: A firestore document with
  - type: URL or file
  - for file a "contents" field
  - a "last refreshed" datetime

## 3. Web Application {:#3.-web}

The web application will provide a user interface for interacting with the
stored instruction snippets.

### 3.1 User Interface (UI) {:#3.1-user}

- **Search Bar**: A prominent search bar allowing users to enter search labels
  or terms.
- **Snippet Display**: A clear and organized list of returned snippets,
  displaying the snippet text and associated labels.
  - each snippet card has snippet text, labels (drawn as pills), a widget/icon
    to copy snippet text to clipboard, and thumbs up and down to add to, and see
    ratings
- **Rating Mechanism**:
  - "Thumbs Up" button for positive feedback.
  - "Thumbs Down" button for negative feedback.
- **Sorting/Filtering Options**:
  - Sort by relevance.
  - Sort by crowdsourced rating (highest rated first).
  - Filter by specific labels.

### 3.2 Functionality {:#3.2-functionality}

- **Snippet Search**: Users can search for snippets using keywords, labels, or a
  combination of both.
- **Snippet Retrieval**: The system will retrieve snippets from Firestore based
  on the user's query.
- **Crowdsourced Rating**: Users can provide feedback on the usefulness of
  snippets through the "thumbs up/down" mechanism. This data will be used to
  calculate a crowdsourced rating for each snippet.
- **Snippet Ranking/Sorting**: Search results can be sorted or ranked based on
  the crowdsourced rating, allowing users to prioritize highly-rated snippets.
- **Admin Panel (Future Consideration)**: A potential future feature for
  administrators to manage labels, review snippets, and monitor system
  performance.
- **User Login**: using firebase identity and google login

## 4. Non-Functional Requirements {:#4.-non-functional}

- **Performance**: The search functionality should be fast and responsive, even
  with a large number of snippets.
- **Scalability**: The system should be designed to scale to accommodate a
  growing number of instructions, snippets, and users.
- **Security**: User data and system access should be secure.
- **Reliability**: The automated processing pipeline should be robust and handle
  errors gracefully.
- **User Experience (UX)**: The web application should be intuitive and easy to
  use.

## 5. Future Enhancements {:#5.-future}

- **User Accounts**: Allow users to save favorite snippets or create
  personalized collections.
- **Snippet Editing**: Enable authorized users to suggest edits or improvements
  to snippets.
- **Advanced Search**: Implement more sophisticated search capabilities, such as
  natural language processing for more nuanced queries.
- **Version Control**: Track changes to snippets over time.
- **Notifications**: Alert users to new relevant snippets based on their
  interests.

