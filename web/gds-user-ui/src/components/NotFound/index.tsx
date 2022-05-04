import { Box, Button, Image, Stack, Text } from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';
import { colors } from 'utils/theme';
import Error404 from 'assets/404-Error.svg';

export default function NotFound() {
  const navigate = useNavigate();
  return (
    <Stack height="100vh" direction="row" justifyContent="center" alignItems="center">
      <Stack>
        <Image src={Error404} width="350px" mx="auto" />
        <Text fontSize="3xl" fontWeight="bold">
          Page not found
        </Text>

        <Text>Sorry, the page you are looking for doesn't exist or has been moved</Text>
        <Box pt={2}>
          <Button onClick={() => navigate('/')} bg={colors.system.blue} color={'white'}>
            Back Home
          </Button>
        </Box>
      </Stack>
    </Stack>
  );
}
