"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Canvas from "./Canvas";
import { Button } from "@mui/material";
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

  return (
    <div
      style={{
        position: "absolute",
        left: "50%",
        top: "50%",
        transform: "translate(-50%, -50%)",
      }}
    >
      <Canvas />
      <Button variant="outlined" onClick={ColorMap} className="mt-5">
        Color
      </Button>
      <Button variant="outlined" onClick={ResetMap} className="mt-5 ml-5">
        Reset
      </Button>
    </div>
  );
}
