import React, { useEffect } from "react";

const Root = () => {
  useEffect(() => {
    if (typeof document !== "undefined") {
      // Safe to use document here
      document.title = "React App";
    }
  }, []);

  return (
    <div>
      <h1>Welcome to React with TypeScript</h1>
    </div>
  );
};

export default Root;
