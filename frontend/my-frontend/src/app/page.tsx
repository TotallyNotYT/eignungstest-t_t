"use client";
import { useCallback, useState } from "react";
import styles from "./page.module.css";

export default function Home() {
  const [myData, setMyData] = useState<string[]>([]);
  const handleClick = useCallback(async () => {
    const data = await fetch("http://localhost:8080/load").then((data) =>
      data.json()
    );
    setMyData(data);
  }, []);
  return (
    <div className={styles.page}>
      <main className={styles.main}>
        hi :D
        <br />
        <div style={{ display: 'flex', gap: "20px"}}>
        <button onClick={handleClick}>click meeee :D</button>
        <button onClick={() => setMyData([])}>clear</button></div>
        {myData.length > 0 && (
          <>
            Meine Daten:
            <br />
            {myData.map((singleItem, index) => {
              return <p key={index}>{singleItem}</p>;
            })}
          </>
        )}
      </main>
    </div>
  );
}
