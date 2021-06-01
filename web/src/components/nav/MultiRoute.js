import { useEffect, useState } from 'react';

const MultiRoute = ({ paths, children }) => {
  const [currentPath, setCurrentPath] = useState(window.location.pathname);

  useEffect(() => {
    const onLocationChange = () => {
      setCurrentPath(window.location.pathname);
    }

    window.addEventListener('popstate', onLocationChange);

    return () => {
      window.removeEventListener('popstate', onLocationChange);
    }
  }, []);

  return paths.has(currentPath) ? children : null;
};

export default MultiRoute;