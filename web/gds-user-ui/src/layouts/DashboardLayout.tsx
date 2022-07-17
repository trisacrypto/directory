import React, { useEffect } from 'react';

import Sidebar from 'components/Sidebar';
import Loader from 'components/Loader';
type DashboardLayoutProp = {
  children: React.ReactNode;
};

const DashboardLayout: React.FC<DashboardLayoutProp> = (props) => {
  return (
    <>
      <Loader />
      {/* <Sidebar {...props} />; */}
    </>
  );
};
export default DashboardLayout;
