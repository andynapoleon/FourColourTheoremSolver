"use client";
import { useEffect, useState, useRef } from "react";
import { useParams, useRouter } from "next/navigation";

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
  const router = useRouter();
  const [mapData, setMapData] = useState<MapData | null>(null);
  const [error, setError] = useState<string>("");
  const [isLoading, setIsLoading] = useState(true);
  const [isDeleting, setIsDeleting] = useState(false);
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

      canvas.width = mapData.width;
      canvas.height = mapData.height;

      const img = new Image();
      img.onload = () => {
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
      };
      img.src = mapData.imageData;
    }
  }, [mapData]);

  const handleDelete = async () => {
    if (
      !window.confirm(
        "Are you sure you want to delete this map? This action cannot be undone."
      )
    ) {
      return;
    }

    const token = localStorage.getItem("token");
    const apiHost = process.env.NEXT_PUBLIC_API_GATEWAY_URL;

    if (!token || !apiHost) {
      setError("Authentication error");
      return;
    }

    setIsDeleting(true);

    try {
      const response = await fetch(`${apiHost}/api/v1/maps/${params.id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      router.push("/profile");
    } catch (err) {
      console.error("Error deleting map:", err);
      setError("Failed to delete map. Please try again later.");
      setIsDeleting(false);
    }
  };

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
          <div className="flex justify-between items-center mb-4">
            <div>
              <h1 className="text-2xl font-bold">{mapData.name}</h1>
              <p className="text-gray-600">
                Created on: {new Date(mapData.createdAt).toLocaleDateString()}
              </p>
            </div>
            <button
              onClick={handleDelete}
              disabled={isDeleting}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg 
                ${
                  isDeleting
                    ? "bg-gray-300 cursor-not-allowed"
                    : "bg-red-500 hover:bg-red-600"
                } 
                text-white transition-colors`}
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
              {isDeleting ? "Deleting..." : "Delete Map"}
            </button>
          </div>
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
