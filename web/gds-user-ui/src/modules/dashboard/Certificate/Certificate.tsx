import BasicDetails from 'components/BasicDetail';
import DashboardLayout from 'layouts/DashboardLayout';
import CertificateLayout from 'layouts/CertificateLayout';

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
