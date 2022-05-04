import { Box, Heading, HStack, Stack } from "@chakra-ui/react";
import NeedsAttention from "components/NeedsAttention";
import NetworkStatus from "components/NetworkStatus";
import OpenSourceResources from "components/OpenSourceResources";
import TravelRuleProviders from "components/TravelRuleProviders";
import TrisaVerifiedLogo from "components/TrisaVerifiedLogo";
import YourImplementation from "components/YourImplementation";

const TechnicalResources: React.FC = () => {
  return (
    <Stack spacing={7}>
      <Heading fontFamily={"'Roboto Slab', serif"}>Technical Resources</Heading>
      <Stack direction={["column", "row"]}>
        <Box width={["100%", "70%"]}>
          <NeedsAttention />
        </Box>
        <Box width={["100%", "30%"]}>
          <NetworkStatus />
        </Box>
      </Stack>
      <YourImplementation />
      <Stack direction={["column", "row"]} spacing={7}>
        <OpenSourceResources />
        <TravelRuleProviders />
      </Stack>
      <TrisaVerifiedLogo />
    </Stack>
  );
};

export default TechnicalResources;
