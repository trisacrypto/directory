import React, { useEffect } from 'react';

import Sidebar from 'components/Sidebar';

type DashboardLayoutProp = {
  children: React.ReactNode;
};

const DashboardLayout: React.FC<DashboardLayoutProp> = (props) => {
  return (
    <>
      <Sidebar {...props} />;
    </>
  );
};
export default DashboardLayout;
