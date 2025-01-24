"use client";
import { useEffect, useState, useRef } from "react";
import { useParams } from "next/navigation";

interface MapData {
  id: string;
  name: string;
  imageData: string;
  matrix: number[][];
  width: number;
  height: number;
  createdAt: string;
}

export default function MapView() {
  const params = useParams();
  const [mapData, setMapData] = useState<MapData | null>(null);
  const [error, setError] = useState<string>("");
  const [isLoading, setIsLoading] = useState(true);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const fetchMapData = async () => {
      const token = localStorage.getItem("token");
      const userId = localStorage.getItem("userId");

      if (!token || !userId) {
        setError("Please login to view this map");
        setIsLoading(false);
        return;
      }

      const apiHost = process.env.NEXT_PUBLIC_API_GATEWAY_URL;
      if (!apiHost) {
        throw new Error("API host is not defined");
      }

      try {
        const response = await fetch(`${apiHost}/api/v1/maps/${params.id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        setMapData(data);
      } catch (err) {
        console.error("Error fetching map:", err);
        setError("Failed to load map. Please try again later.");
      } finally {
        setIsLoading(false);
      }
    };

    fetchMapData();
  }, [params.id]);

  useEffect(() => {
    if (mapData && canvasRef.current) {
      const canvas = canvasRef.current;
      const ctx = canvas.getContext("2d");
      if (!ctx) return;

      // Set canvas dimensions
      canvas.width = mapData.width;
      canvas.height = mapData.height;

      // Load and draw the base image
      const img = new Image();
      img.onload = () => {
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
      };
      img.src = mapData.imageData;
    }
  }, [mapData]);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      </div>
    );
  }

  if (!mapData) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-gray-600">Map not found</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-lg shadow-lg p-6">
          <h1 className="text-2xl font-bold mb-4">{mapData.name}</h1>
          <p className="text-gray-600 mb-4">
            Created on: {new Date(mapData.createdAt).toLocaleDateString()}
          </p>
          <div className="relative w-full overflow-auto">
            <canvas
              ref={canvasRef}
              className="border border-gray-200 rounded mx-auto"
              style={{
                maxWidth: "100%",
                height: "auto",
              }}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
