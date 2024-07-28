"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Canvas from "./Canvas";
import { ColorMap, ResetMap } from "./Canvas";

export default function HomeContent() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/login");
    } else {
      setIsLoading(false);
    }
  }, [router]);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  const buttonStyles =
    "px-4 py-2 text-sm font-bold uppercase rounded-full shadow-md transition-all duration-300 ease-in-out transform hover:-translate-y-1 hover:shadow-lg active:translate-y-0 active:shadow-md";

  return (
    <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2">
      <Canvas />
      <div className="mt-5 flex justify-center space-x-4">
        <button
          onClick={ColorMap}
          className={`${buttonStyles} bg-green-500 text-white border-2 border-green-600 hover:bg-green-600`}
        >
          Color
        </button>
        <button
          onClick={ResetMap}
          className={`${buttonStyles} bg-red-500 text-white border-2 border-red-600 hover:bg-red-600`}
        >
          Reset
        </button>
      </div>
    </div>
  );
}
