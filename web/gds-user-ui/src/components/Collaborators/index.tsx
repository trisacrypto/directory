import { Heading, Stack } from "@chakra-ui/react";
import CollaboratorsSection from "components/CollaboratorsSection";

const Collaborators: React.FC = () => {
  return (
    <Stack>
      <Heading marginBottom="69px">Collaborators</Heading>
      <CollaboratorsSection />
    </Stack>
  );
};

export default Collaborators;
