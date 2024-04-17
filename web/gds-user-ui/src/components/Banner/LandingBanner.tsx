import { Alert, AlertDescription, AlertIcon, Link, Stack } from "@chakra-ui/react";
import { Trans } from "@lingui/macro";

const LandingBanner = () => {
  return (
    <Stack spacing={3}>
      <Alert
      status="info"
      justifyContent="center"
      paddingY={'5'}
      fontSize={'lg'}
      >
      <AlertIcon />
      <AlertDescription>
        <Link href="https://calendar.app.google/FBg7GTmgDfeMbUMT9" isExternal>
          <Trans id="Schedule a demo to learn about TRISA's open source self-hosted solution for cost-effective Travel Rule compliance.">
            Schedule a demo to learn about TRISA's open source self-hosted solution for cost-effective Travel Rule compliance.
          </Trans>
        </Link>
      </AlertDescription>
      </Alert>
    </Stack>
  );
};

export default LandingBanner;
