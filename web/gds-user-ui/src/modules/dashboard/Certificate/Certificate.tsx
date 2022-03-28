import { SimpleDashboardLayout } from 'layouts';
import CertificateLayout from 'layouts/CertificateLayout';

const Certificate: React.FC = () => {
  return (
    // <DashboardLayout>
    //   <CertificateLayout>
    //     <BasicDetails />
    //   </CertificateLayout>
    // </DashboardLayout>
    <SimpleDashboardLayout>
      <CertificateLayout />
    </SimpleDashboardLayout>
  );
};

export default Certificate;
