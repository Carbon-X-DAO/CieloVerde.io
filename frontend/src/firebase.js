import firebase from 'firebase/compat/app'
import "firebase/compat/database"

const firebaseConfig = {
    apiKey: "AIzaSyD0nAAuPxwUXsoZ3l1m_U2hEQxqJ3BfVyM",
    authDomain: "contact-form-5c442.firebaseapp.com",
    projectId: "contact-form-5c442",
    storageBucket: "contact-form-5c442.appspot.com",
    messagingSenderId: "896627853267",
    appId: "1:896627853267:web:44c40427518911c84e9126",
    measurementId: "G-XXFLSCZZSB"
  };

  const fireDB = firebase.initializeApp(firebaseConfig)

  export default fireDB.database().ref();