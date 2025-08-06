import { initializeApp, type FirebaseOptions } from 'firebase/app';
import { getFirestore } from 'firebase/firestore';

const firebaseConfig: FirebaseOptions = {
	apiKey: 'AIzaSyDF-r3ZvAiZz32_ITugrRFDERJaB4NF45s',
	authDomain: 'new-test-297222.firebaseapp.com',
	projectId: 'new-test-297222',
	storageBucket: 'new-test-297222.appspot.com',
	messagingSenderId: '298870980814',
	appId: '1:298870980814:web:235ce1376aaf86b99b909f',
	measurementId: 'G-9ER4KVQE61'
};

const app = initializeApp(firebaseConfig);
export const db = getFirestore(app);
