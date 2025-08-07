const admin = require('firebase-admin');

// Initialize the admin SDK
admin.initializeApp({
  credential: admin.credential.applicationDefault(),
  projectId: 'new-test-297222',
});

const db = admin.firestore();



async function deleteQueryBatch(query, resolve, reject) {
  try {
    const snapshot = await query.get();

    if (snapshot.size === 0) {
      return resolve();
    }

    const batch = db.batch();
    snapshot.docs.forEach((doc) => {
      batch.delete(doc.ref);
    });

    await batch.commit();

    // Recurse on the next process tick, to avoid
    // exploding the stack.
    process.nextTick(() => {
      deleteQueryBatch(query, resolve, reject);
    });
  } catch (err) {
    reject(err);
  }
}

async function main() {
  console.log('Deleting snippets collection...');
    const snippetsQuery = db.collection('snippets').orderBy('__name__').limit(20);
  await new Promise((resolve, reject) => deleteQueryBatch(snippetsQuery, resolve, reject));
  console.log('Snippets collection deleted.');

  console.log('Deleting sources collection...');
    const sourcesQuery = db.collection('sources').orderBy('__name__').limit(20);
  await new Promise((resolve, reject) => deleteQueryBatch(sourcesQuery, resolve, reject));
  console.log('Sources collection deleted.');
}

main().catch(console.error);
