const { initializeApp } = require('firebase/app');
const { getFirestore, collection, addDoc, Timestamp } = require('firebase/firestore');

const firebaseConfig = {
  apiKey: "AIzaSyDF-r3ZvAiZz32_ITugrRFDERJaB4NF45s",
  authDomain: "new-test-297222.firebaseapp.com",
  projectId: "new-test-297222",
  storageBucket: "new-test-297222.appspot.com",
  messagingSenderId: "298870980814",
  appId: "1:298870980814:web:235ce1376aaf86b99b909f",
  measurementId: "G-9ER4KVQE61"
};

const app = initializeApp(firebaseConfig);
const db = getFirestore(app);

const snippets = [
  {
    text: "When writing Python code, always adhere to the PEP 8 style guide. This includes using 4 spaces for indentation, limiting lines to 79 characters, and using snake_case for variable and function names. Consistent styling makes code more readable and maintainable.",
    labels: ["python", "pep8", "style-guide", "coding-conventions"],
  },
  {
    text: "For JavaScript projects, use Prettier to automatically format your code. This eliminates debates about style and ensures consistency across the entire codebase. Configure it to run on save and as a pre-commit hook.",
    labels: ["javascript", "prettier", "code-formatting", "tooling"],
  },
  {
    text: "In React, prefer functional components and hooks over class-based components. Functional components are more concise, easier to test, and align better with the modern React paradigm. Use the `useState` and `useEffect` hooks for state and side effects.",
    labels: ["react", "javascript", "functional-components", "hooks"],
  },
  {
    text: "When designing a REST API, use clear and consistent naming conventions for your endpoints. Use plural nouns for resource collections (e.g., `/users`, `/products`) and standard HTTP methods (GET, POST, PUT, DELETE) for operations.",
    labels: ["rest-api", "api-design", "best-practices"],
  },
  {
    text: "Always include a `.gitignore` file in your Git repositories. This file tells Git which files and directories to ignore, such as `node_modules`, build artifacts, and environment-specific files. This keeps your repository clean and focused on the source code.",
    labels: ["git", "version-control", "gitignore"],
  },
];

async function populateDatabase() {
  const snippetsCollection = collection(db, 'snippets');
  for (const snippet of snippets) {
    await addDoc(snippetsCollection, {
      ...snippet,
      created_at: Timestamp.now(),
      thumbs_up: 0,
      thumbs_down: 0,
    });
  }
  console.log('Database populated successfully!');
  process.exit(0);
}

populateDatabase();
