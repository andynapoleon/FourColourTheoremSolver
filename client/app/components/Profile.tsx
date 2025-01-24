"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import styles from "./styles/Profile.module.css";

interface Map {
  id: string;
  name: string;
  createdAt: string;
}

export default function Profile() {
  const router = useRouter();
  const [userName, setUserName] = useState("");
  const [userEmail, setUserEmail] = useState("");
  const [maps, setMaps] = useState<Map[] | null>(null); // Initialize as null
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchUserData = async () => {
      const token = localStorage.getItem("token");
      const userId = localStorage.getItem("userId");

      if (!token || !userId) {
        router.push("/login");
        return;
      }

      const apiHost = process.env.NEXT_PUBLIC_API_GATEWAY_URL;
      if (!apiHost) {
        throw new Error("API host is not defined");
      }

      const storedUserName = localStorage.getItem("name");
      setUserName(storedUserName || "User");

      try {
        const response = await fetch(
          `${apiHost}/api/v1/maps?userId=${userId}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        setMaps(Array.isArray(data) ? data : []);
      } catch (err) {
        console.error("Error fetching maps:", err);
        setError("Failed to load maps. Please try again later.");
        setMaps([]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchUserData();
  }, [router]);

  const handleMapClick = (mapId: string) => {
    router.push(`/maps/${mapId}`);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className={styles.profileContainer}>
      <div className={styles.profileCard}>
        <h1 className={styles.title}>Profile</h1>

        <div className={styles.profileSection}>
          <h2>Personal Information</h2>
          <div className={styles.infoGroup}>
            <label>Name:</label>
            <p>{userName}</p>
          </div>
          <div className={styles.infoGroup}>
            <label>Email:</label>
            <p>{userEmail || "email@example.com"}</p>
          </div>
        </div>

        <div className={styles.profileSection}>
          <h2>Maps</h2>
          {error && (
            <div className="text-red-500 text-center my-4">{error}</div>
          )}
          {!error && maps && maps.length === 0 ? (
            <p className="text-gray-500 text-center my-4">No maps found</p>
          ) : (
            <div className={styles.mapsGrid}>
              {maps?.map((map) => (
                <div
                  key={map.id}
                  className={styles.mapCard}
                  onClick={() => handleMapClick(map.id)}
                >
                  <h3>{map.name}</h3>
                  <p className={styles.mapDate}>
                    Created: {new Date(map.createdAt).toLocaleDateString()}
                  </p>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className={styles.profileSection}>
          <h2>Account Settings</h2>
          <button
            className={styles.button}
            onClick={() => alert("Feature coming soon (or not)!")}
          >
            Change Password
          </button>
        </div>
      </div>
    </div>
  );
}
