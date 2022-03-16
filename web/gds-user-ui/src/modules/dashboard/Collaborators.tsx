import { Heading, Stack } from '@chakra-ui/react';
import CollaboratorsSection from 'components/CollaboratorsSection';
import DashboardLayout from 'layouts/DashboardLayout';

const Collaborators: React.FC = () => {
  return (
    <DashboardLayout>
      <Heading marginBottom="69px">Collaborators</Heading>
      <CollaboratorsSection />
    </DashboardLayout>
  );
};

export default Collaborators;
