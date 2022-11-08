import { Heading } from '@chakra-ui/react';
import CollaboratorsSection from 'components/CollaboratorsSection';
// import DashboardLayout from 'layouts/DashboardLayout';

const Collaborators: React.FC = () => {
  return (
    <>
      <Heading marginBottom="69px">Collaborators</Heading>
      <CollaboratorsSection />
    </>
  );
};

export default Collaborators;
