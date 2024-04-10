import { Alert, AlertDescription, AlertIcon } from "@chakra-ui/react";

const LandingBanner = () => {
  return (
    <Alert
    status="info"
    variant="subtle"
    flexDirection="column"
    alignItems="center"
    justifyContent="center"
    textAlign="center"
    height="150px"
  >
    <AlertIcon boxSize="30px" my={1} mr={0} />
    <AlertDescription mt={4} maxWidth="sm" fontWeight={'semibold'}>
      {/* Add link to schedule a demo */}
    Schedule a demo to learn about TRISA's open source self-hosted solution for cost-effective Travel Rule compliance.
    </AlertDescription>
  </Alert>
  );
};

export default LandingBanner;
