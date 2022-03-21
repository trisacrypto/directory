import BasicDetails from 'components/BasicDetail';
import CertificateLayout from 'layouts/CertificateLayout';
import DashboardLayout from 'layouts/DashboardLayout';

const Certificate: React.FC = () => {
  return (
    <DashboardLayout>
      <CertificateLayout>
        <BasicDetails />
      </CertificateLayout>
    </DashboardLayout>
  );
};

export default Certificate;
