"use client";

import React, { useEffect, useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import styles from "./styles/ProfileContent.module.css";

interface UserData {
  name: string;
  email: string;
  joinDate: string;
}

interface SavedImage {
  id: string;
  url: string;
  title: string;
}

const ProfileContent: React.FC = () => {
  const [userData, setUserData] = useState<UserData | null>(null);
  const [savedImages, setSavedImages] = useState<SavedImage[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        setIsLoading(true);
        const token = localStorage.getItem("token");
        if (!token) {
          router.push("/login");
          return;
        }

        const response = await fetch("/api/user-profile", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        if (!response.ok) {
          throw new Error("Failed to fetch user data");
        }
        const data = await response.json();
        setUserData(data.userData);
        setSavedImages(data.savedImages);
      } catch (err) {
        setError("An error occurred while fetching user data");
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchUserData();
  }, [router]);

  if (isLoading) {
    return <div className={styles.centerMessage}>Loading...</div>;
  }

  if (error) {
    return <div className={styles.centerMessage}>Error: {error}</div>;
  }

  return (
    <div className={styles.profileContainer}>
      {userData && (
        <div className={styles.userInfo}>
          <h2>{userData.name}</h2>
          <p>Email: {userData.email}</p>
          <p>Joined: {new Date(userData.joinDate).toLocaleDateString()}</p>
        </div>
      )}

      <h2>Saved Images</h2>
      {savedImages.length > 0 ? (
        <div className={styles.imageGrid}>
          {savedImages.map((image) => (
            <div key={image.id} className={styles.imageCard}>
              <Image
                src={image.url}
                alt={image.title}
                width={200}
                height={200}
                layout="responsive"
              />
              <p>{image.title}</p>
            </div>
          ))}
        </div>
      ) : (
        <p>No saved images yet.</p>
      )}
    </div>
  );
};

export default ProfileContent;
