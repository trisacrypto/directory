import { Box, Button, Image, Stack, Text, Link } from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';
import { colors } from 'utils/theme';
import Error404 from 'assets/404-Error.svg';
import TrisaLogo from 'assets/trisa-logo-white.png';
import MaintenanceSVG from 'assets/maintenance.svg';
import { Trans } from '@lingui/react';

const Maintenance: React.FC = () => {
  return (
    <Stack direction="row" justifyContent="center" alignItems="center" textAlign={'center'}>
      <Stack fontSize="xl" pt={'80px'} mx={{ lg: 365 }}>
        <Image src={TrisaLogo} mx="auto" width={64} />
        <Text fontSize="3xl">
          <Trans id="We’ll be back soon.">We’ll be back soon.</Trans>
        </Text>
        <Text>
          <Trans id="The TRISA Global Directory Service (GDS) is temporarily undergoing maintenance. Please try again later.">
            The TRISA Global Directory Service (GDS) is temporarily undergoing maintenance. Please
            try again later.
          </Trans>
        </Text>
        <Stack alignItems={'center'} mx={'auto'}>
          <Image src={MaintenanceSVG} width={'50%'} pt={8} loading="eager" />
        </Stack>

        <Text pt={'5px'}>
          <Trans id="Join">Join</Trans>{' '}
          <Link
            isExternal
            textDecoration={'underline'}
            color={'#1F4CED'}
            href="https://trisa-workspace.slack.com/">
            <Trans id="TRISA's Slack channel">TRISA's Slack channel</Trans>
          </Link>{' '}
          <Trans id="to receive maintenance and outage notifications. If you have an immediate concern, please email support@rotational.io.">
            to receive maintenance and outage notifications. If you have an immediate concern,
            please email support@rotational.io.
          </Trans>
        </Text>
      </Stack>
    </Stack>
  );
};

export default Maintenance;