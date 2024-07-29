"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Canvas from "./Canvas";
import { handleColorMap, handleResetMap, handleDownloadMap } from "./Canvas";

export default function HomeContent() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(true);
  const [userName, setUserName] = useState("");

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/login");
    } else {
      setIsLoading(false);
      // Retrieve user name from localStorage
      const storedUserName = localStorage.getItem("name");
      setUserName(storedUserName || "User"); // Default to "User" if name not found
    }
  }, [router]);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  const buttonStyles =
    "px-4 py-2 text-sm font-bold uppercase rounded-full shadow-md transition-all duration-300 ease-in-out transform hover:-translate-y-1 hover:shadow-lg active:translate-y-0 active:shadow-md";

  return (
    <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2">
      <h1 className="text-2xl font-bold text-center mb-6 text-black-600">
        <span>Welcome, </span> {userName}!
      </h1>
      <p className="text-center mb-6">Please draw a map using your mouse 😊!</p>
      <Canvas />
      <div className="mt-5 flex justify-center space-x-4">
        <button
          onClick={handleColorMap}
          className={`${buttonStyles} bg-green-500 text-white border-2 border-green-600 hover:bg-green-600`}
        >
          Color
        </button>
        <button
          onClick={handleResetMap}
          className={`${buttonStyles} bg-red-500 text-white border-2 border-red-600 hover:bg-red-600`}
        >
          Reset
        </button>
        <button
          onClick={handleDownloadMap}
          className={`${buttonStyles} bg-blue-500 text-white border-2 border-blue-600 hover:bg-blue-600`}
        >
          Download
        </button>
      </div>
    </div>
  );
}
