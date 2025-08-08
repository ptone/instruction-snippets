# Project Milestones

This document outlines the phased development milestones for the LLM Prompt Engineering Instruction Management System.

## Milestone 1: Core Snippet Display and Frontend Setup

**Goal:** Establish the basic frontend architecture and display manually-added snippets from Firestore. This milestone validates the core technology choices and data structure.

*   **Backend:**
    *   Define and create the Firestore database structure for `snippets` and `sources`.
    *   Manually populate the `snippets` collection with 3-5 sample documents. Note this may likely be done from a prior iteration, with sample data from @scripts/populate.js

*   **Frontend:**
    *   An initialized frontend project using sveltekit is in the frontend folder
    *   Set up Firebase configuration for web client.
    *   Create a main view to fetch and display all snippets from the Firestore `snippets` collection. 
    *   For each snippet, display the text in a card component and associated labels as pills.

    The Firebase config is:

    const firebaseConfig = {
        apiKey: "AIzaSyDF-r3ZvAiZz32_ITugrRFDERJaB4NF45s",
        authDomain: "new-test-297222.firebaseapp.com",
        projectId: "new-test-297222",
        storageBucket: "new-test-297222.firebasestorage.app",
        messagingSenderId: "298870980814",
        appId: "1:298870980814:web:235ce1376aaf86b99b909f",
        measurementId: "G-9ER4KVQE61"
    };

## Milestone 2: Snippet Search and Rating Functionality

**Goal:** Implement the primary user interaction features: searching for snippets and rating them.

*   **Frontend:**
    *   Add a prominent search bar to the UI.
    *   Implement client-side logic to filter displayed snippets based on user input (searching against snippet text and labels).
    *   Add "Thumbs Up" and "Thumbs Down" icons/buttons to each snippet card.
    *   Implement the functionality to update the corresponding `thumbs_up` or `thumbs_down` counter in Firestore when a user clicks a rating button.
    *   Display the current rating counts on each snippet.
    *   Add a "copy to clipboard" button for the snippet text.

## Milestone 3: Automated Backend Processing Pipeline

**Goal:** Build the core backend service that can ingest a markdown document, process it using an LLM, and store the resulting snippets in Firestore.

*   **Backend (Go):**
    *   Set up a basic Go server project, deployable to Cloud Run.
    *   Create an HTTP endpoint or a manually-triggered function that accepts markdown content.
    *   Integrate with the Vertex AI Gemini API.
    *   Implement the processing logic:
        1. Store the markdown source in a `sources` collection in Firestore, noting the generated document id
        1.  Send the markdown content to the LLM to be broken down into discrete instruction snippets.
        1.  For each generated snippet, call the LLM again to generate relevant labels.
        1.  For each snippet, call the Vertex AI API for an embedding model to generate an embedding vector.
    *   Save the processed snippets, labels, and embeddings to the Firestore `snippets` collection, include a doc-reference to the sources document from which the snippet is associated
    *   Note the "last processed" timestamp on the sources document

## Milestone 3.5: Make the ingest more resilient

As part of processing a source file, we need some sort of key. For now we will use file-name or URL. We will still use the document id as an internal link for relationships as the we may involve to a more complex duplicate detecting key.

When a source file is being submitted for re-processing, the order should be:

        1. Delete all snippets associated with that source
        1. update the contents of the source from the provided content
        1. reprocess the source to generate snippets

This milestone should support the ability to re-run integration tests with the sample files or URLs without producing a proliferation of firestore documents

## Milestone 4: End-to-End Flow with User Authentication and Advanced Features

**Goal:** Create a full end-to-end experience, from automated ingestion to a feature-rich user interface.

*   **Backend:**

    *   Implement the Vertex AI safety filter for all extracted snippet content. Start by simply storing safety score with snippet for client side filter.
    *   Ensure the URL is used as a source key for URL based sources

*   **Frontend:**
    *   Implement user login using Firebase Authentication with Google as the provider.
    *   Add sorting options to the search results (e.g., sort by rating, sort by date).
    *   Add advanced filtering capabilities (e.g., filter by one or more specific label pills).
    *   Add a "add contribution" feature, that will allow users to submit a form with either file upload or by URL (but not both, use selector in form). Style everything with UI components

## Milestone 5: Refinement and Advanced Search

**Goal:** Polish the application, improve search relevance, and prepare for future enhancements.

*   **Backend:**
    *   Develop a backend endpoint to support semantic search. This endpoint will take a user's query, generate an embedding for it, and use it to find the most similar snippet embeddings in Firestore (Vector Search).

*   **Frontend:**
    *   Integrate the frontend search bar with the new semantic search backend endpoint for more relevant search results.
    *   Refine the overall UI/UX for a polished, intuitive, and professional feel.
    *   Review and implement key items from the "Future Enhancements" section of the PRD, such as allowing users to save favorite snippets.


# TODO collector

    *   Update the mechanism for ingesting content from URLs,  triggered by a Cloud Scheduler job
        * This should check when it was last refreshed, and perhaps retrieve and store a last-modified date from server or ETag to check
        * Once the candidates are identified for updating, they should be enqued for processing. This could be done in a set of go-routines, ulimately as part of a Cloud Run Job (future work)