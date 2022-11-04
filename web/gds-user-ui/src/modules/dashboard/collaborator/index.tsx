import { useEffect, useCallback } from 'react';
import { Heading } from '@chakra-ui/react';
import CollaboratorsSection from 'components/CollaboratorsSection';
// import DashboardLayout from 'layouts/DashboardLayout';
import { getAllCollaborators } from './service';

const Collaborators: React.FC = () => {
  const fetchAllCollaborators = useCallback(async () => {
    try {
      const res = await getAllCollaborators();
      console.log(res);
    } catch (err) {
      console.log(err);
    }
  }, []);

  useEffect(() => {
    fetchAllCollaborators();
  }, [fetchAllCollaborators]);

  return (
    <>
      <Heading marginBottom="69px">Collaborators</Heading>
      <CollaboratorsSection />
    </>
  );
};

export default Collaborators;
