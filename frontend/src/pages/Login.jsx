import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import loginImage from "../assets/images/login/login.png";
export default function UnifiedAuthPage() {
  const [step, setStep] = useState("signin"); // signin, signup, otp
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [otp, setOtp] = useState(["", "", "", "", "", ""]);
  const [resendTimer, setResendTimer] = useState(0);
  const [canResend, setCanResend] = useState(false);
  const navigate = useNavigate();

  // Timer countdown effect
  useEffect(() => {
    let interval;
    if (resendTimer > 0) {
      interval = setInterval(() => {
        setResendTimer((prev) => {
          if (prev <= 1) {
            setCanResend(true);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [resendTimer]);

  // Form input handler
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  // OTP input handler
  const handleOtpChange = (e, index) => {
    const value = e.target.value;
    if (!/^[0-9]?$/.test(value)) return;

    const newOtp = [...otp];
    newOtp[index] = value;
    setOtp(newOtp);

    if (value && index < otp.length - 1) {
      document.getElementById(`otp-${index + 1}`).focus();
    }
  };

  // OTP backspace handler
  const handleKeyDown = (e, index) => {
    if (e.key === "Backspace" && !otp[index] && index > 0) {
      document.getElementById(`otp-${index - 1}`).focus();
    }
  };

  // Sign In
  const handleSignIn = async () => {
    if (!formData.email || !formData.password) {
      alert("Please fill in all required fields.");
      return;
    }

    try {
      const response = await axios.post("http://localhost:9000/login", {
        email: formData.email,
        password: formData.password,
      });

      console.log("Login successful:", response.data);
      alert("Login successful!");
      navigate("/");
    } catch (error) {
      console.error("Login error:", error);
      alert(
        error.response?.data?.message ||
          "Login failed. Please check your credentials."
      );
    }
  };

  // Sign Up
  const handleSignUp = async () => {
    if (!formData.email || !formData.password) {
      alert("Please fill in all required fields.");
      return;
    }

    try {
      const response = await axios.post("http://localhost:9000/signup", {
        email: formData.email,
        password: formData.password,
      });

      console.log("Signup successful:", response.data);
      alert("Account created! Please verify OTP.");
      setStep("otp");
      setResendTimer(60); // Start 60 second timer
      setCanResend(false);
    } catch (error) {
      console.error("Signup error:", error);
      alert(
        error.response?.data?.message || "Signup failed. Please try again."
      );
    }
  };

  // Verify OTP
  const handleVerifyOtp = async () => {
    const fullOtp = otp.join("");

    if (fullOtp.length !== 6) {
      alert("Please enter all 6 digits.");
      return;
    }

    try {
      const response = await axios.post("http://localhost:9000/verify-otp", {
        email: formData.email,
        otp: fullOtp,
      });

      console.log("OTP verified:", response.data);
      alert("Verification successful!");
      navigate("/");
    } catch (error) {
      console.error("OTP verification error:", error);
      alert(
        error.response?.data?.message ||
          "OTP verification failed. Please try again."
      );
    }
  };

  // Resend OTP
  const handleResendOtp = async () => {
    if (!canResend && resendTimer > 0) return;

    try {
      await axios.post("http://localhost:9000/resend-otp", {
        email: formData.email,
      });
      alert("OTP resent successfully!");
      setOtp(["", "", "", "", "", ""]);
      setResendTimer(60); // Reset timer to 60 seconds
      setCanResend(false);
      // Focus first OTP input
      document.getElementById("otp-0")?.focus();
    } catch (error) {
      console.error("Resend OTP error:", error);
      alert("Failed to resend OTP. Please try again.");
    }
  };

  // Format timer display
  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, "0")}`;
  };

  // Render Sign In Form
  const renderSignIn = () => (
    <>
      <h1 className="text-4xl font-bold text-blue-700 mb-4 text-center">
        Sign In
      </h1>
      <p className="text-gray-600 text-center mb-8">
        Welcome back! Please enter your credentials to access your account.
      </p>

      <div className="space-y-5">
        <input
          type="email"
          name="email"
          placeholder="Enter your email"
          value={formData.email}
          onChange={handleChange}
          className="w-full px-4 py-3 bg-white rounded-lg border border-blue-100 focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
        />

        <input
          type="password"
          name="password"
          placeholder="Password"
          value={formData.password}
          onChange={handleChange}
          className="w-full px-4 py-3 bg-white rounded-lg border border-blue-100 focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
        />

        <div className="grid grid-cols-2 gap-4 pt-4">
          <button
            onClick={handleSignIn}
            className="px-6 py-3 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition shadow-md"
          >
            Login
          </button>

          <button
            onClick={() => {
              setStep("signup");
              setFormData({ email: "", password: "" });
            }}
            className="px-6 py-3 bg-white text-blue-600 font-semibold rounded-lg border-2 border-blue-600 hover:bg-blue-50 transition"
          >
            Sign Up
          </button>
        </div>
      </div>
    </>
  );

  // Render Sign Up Form
  const renderSignUp = () => (
    <>
      <h1 className="text-4xl font-bold text-blue-700 mb-8 text-center">
        Create Account
      </h1>

      <div className="space-y-5">
        <input
          type="email"
          name="email"
          placeholder="Enter your email"
          value={formData.email}
          onChange={handleChange}
          className="w-full px-4 py-3 bg-white rounded-lg border border-blue-100 focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
        />

        <input
          type="password"
          name="password"
          placeholder="Password"
          value={formData.password}
          onChange={handleChange}
          className="w-full px-4 py-3 bg-white rounded-lg border border-blue-100 focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
        />

        <div className="grid grid-cols-2 gap-4 pt-4">
          <button
            onClick={handleSignUp}
            className="px-6 py-3 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition shadow-md"
          >
            Sign Up
          </button>
          <button
            onClick={() => {
              setStep("signin");
              setFormData({ email: "", password: "" });
            }}
            className="px-6 py-3 bg-white text-blue-600 font-semibold rounded-lg border-2 border-blue-600 hover:bg-blue-50 transition"
          >
            Sign In
          </button>
        </div>
      </div>
    </>
  );

  // Render OTP Form
  const renderOtp = () => (
    <>
      <h1 className="text-4xl font-bold text-blue-700 mb-6 text-center">
        Verify OTP
      </h1>
      <p className="text-gray-600 text-center mb-8">
        Enter the 6-digit code sent to <br />
        <span className="font-semibold text-blue-700">{formData.email}</span>
      </p>

      <div className="flex justify-center gap-3 mb-8">
        {otp.map((digit, index) => (
          <input
            key={index}
            id={`otp-${index}`}
            type="text"
            maxLength="1"
            value={digit}
            onChange={(e) => handleOtpChange(e, index)}
            onKeyDown={(e) => handleKeyDown(e, index)}
            className="w-12 h-12 text-center text-2xl font-semibold text-blue-700 bg-white border border-blue-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
          />
        ))}
      </div>

      <button
        onClick={handleVerifyOtp}
        className="w-full px-6 py-3 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition shadow-md mb-6"
      >
        Verify OTP
      </button>

      {/* Resend OTP Section */}
      <div className="bg-white bg-opacity-50 rounded-lg p-4 mb-4">
        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-600">
            {resendTimer > 0 ? (
              <>
                Resend OTP in{" "}
                <span className="font-semibold text-blue-700">
                  {formatTime(resendTimer)}
                </span>
              </>
            ) : (
              "Didn't receive the code?"
            )}
          </span>
          <button
            onClick={handleResendOtp}
            disabled={!canResend && resendTimer > 0}
            className={`text-sm font-semibold px-4 py-2 rounded-lg transition ${
              canResend || resendTimer === 0
                ? "text-blue-600 hover:bg-blue-100 cursor-pointer"
                : "text-gray-400 cursor-not-allowed"
            }`}
          >
            Resend OTP
          </button>
        </div>
      </div>

      <p className="text-center text-sm text-gray-600">
        Wrong email?{" "}
        <span
          onClick={() => {
            setStep("signup");
            setOtp(["", "", "", "", "", ""]);
            setResendTimer(0);
            setCanResend(false);
          }}
          className="text-blue-600 font-medium cursor-pointer hover:underline"
        >
          Go back
        </span>
      </p>
    </>
  );

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <div className="flex flex-col lg:flex-row items-center justify-between w-full max-w-7xl gap-8">
        {/* Left Image Section */}
          <div className="flex justify-center items-center w-1/2 -ml-20">
            <img
              src={loginImage}
              alt="Authentication Illustration"
              className="w-[650px] h-[650px] object-cover"
            />
          </div>


        {/* Right Form Section */}
        <div className="w-full lg:w-[500px] bg-gradient-to-br from-white via-blue-50 to-blue-200 rounded-2xl shadow-2xl p-8 lg:p-12 flex flex-col justify-center">
          {step === "signin" && renderSignIn()}
          {step === "signup" && renderSignUp()}
          {step === "otp" && renderOtp()}
        </div>
      </div>
    </div>
  );
}
