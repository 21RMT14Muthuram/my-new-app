import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import image from "../assets/images/login/login.png";

export default function SignInPage() {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });

  const navigate = useNavigate();

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async () => {
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

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="flex flex-row items-center justify-between w-full max-w-7xl">
        {/* Left Image Section */}
        <div className="flex justify-center items-center w-1/2 -ml-20">
          <img
            src={image}
            alt="Sign In Illustration"
            className="w-[650px] h-[650px] object-cover"
          />
        </div>

        {/* Right Form Section */}
        <div className="w-[500px] bg-gradient-to-br from-white via-blue-50 to-blue-200 rounded-2xl shadow-2xl p-12 flex flex-col justify-center ml-10">
          <h1 className="text-4xl font-bold text-blue-700 mb-8 text-center">
            Sign In
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
                onClick={handleSubmit}
                className="px-6 py-3 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition shadow-md"
              >
                Login
              </button>

              <button
                onClick={() => navigate("/login")}
                className="px-6 py-3 bg-white text-blue-600 font-semibold rounded-lg border-2 border-blue-600 hover:bg-blue-50 transition"
              >
                Sign Up
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
