import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import SignUpPage from "./pages/login";
import OtpForm from "./pages/Otp";
import SignInPage from "./pages/SignInPage";
import HomePage from "./pages/HomePage";
import UnifiedAuthPage from "./pages/login";


function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<UnifiedAuthPage />} />
        <Route path="/otp" element={<OtpForm />} />
        <Route path="/signin" element={<SignInPage />} />
      </Routes>
    </Router>
  );
}

export default App;
