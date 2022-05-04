import Sidebar from 'components/Sidebar';

type DashboardLayoutProp = {
  children: React.ReactNode;
};

const DashboardLayout: React.FC<DashboardLayoutProp> = (props) => <Sidebar {...props} />;

export default DashboardLayout;
