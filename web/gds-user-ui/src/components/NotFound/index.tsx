import { Box, Button, Stack, Text } from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';
import { colors } from 'utils/theme';
import Error404 from 'assets/404-Error.svg';
import { Trans } from '@lingui/react';
import CkLazyLoadImage from 'components/LazyImage';

export default function NotFound() {
  const navigate = useNavigate();
  return (
    <Stack direction="row" justifyContent="center" alignItems="center" textAlign={'center'}>
      <Stack fontSize="2xl">
        <CkLazyLoadImage src={Error404} width="350px" mx="auto" />
        <Text fontSize="4xl" fontWeight="bold">
          <Trans id="Page Not Found">Page Not Found</Trans>
        </Text>

        <Text>
          <Trans id="Sorry, the page you are looking for doesn't exist or has been moved">
            Sorry, the page you are looking for doesn't exist or has been moved
          </Trans>
        </Text>
        <Box pt={2}>
          <Button onClick={() => navigate('/')} bg={colors.system.blue} color={'white'}>
            <Trans id="Back Home">Back Home</Trans>
          </Button>
        </Box>
      </Stack>
    </Stack>
  );
}
