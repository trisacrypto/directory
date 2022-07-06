import React, { useEffect } from 'react';
import useAuth from 'hooks/useAuth';
import Sidebar from 'components/Sidebar';

type DashboardLayoutProp = {
  children: React.ReactNode;
};

const DashboardLayout: React.FC<DashboardLayoutProp> = (props) => {
  console.log('[DashboardLayout]');

  return (
    <>
      <Sidebar {...props} />;
    </>
  );
};
export default DashboardLayout;
