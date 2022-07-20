import { Box, Button, Stack, Text } from '@chakra-ui/react';
import Card from 'components/Card/Card';
import { BsCheckCircleFill } from 'react-icons/bs';
import TrisaLogo from 'components/Icon/TrisaLogo';

const TrisaVerifiedLogo: React.FC = () => {
  return (
    <Card maxW="100%">
      <Card.CardHeader>TRISA Verified Logo</Card.CardHeader>
      <Card.CardBody>
        <Stack direction={['column', 'row']} spacing={5}>
          <Box maxW={['100%', '50%']}>
            <Text>
              TRISA verified members may download and display a “TRISA Verified VASP” logo on their
              website. The logo is unique to your VASP and non-reproducible. Members may download
              their logo after verification is complete and their certificate has been issued. The
              logo is in .svg fotmat
            </Text>
          </Box>
          <Box display="flex" justifyContent="center" width="100%">
            <Stack spacing={2} alignItems="center">
              <Box
                justifyContent="center"
                background="#E5EDF1"
                border="1px solid #23A7E0"
                borderRadius={10}
                position="relative"
                paddingX="43px">
                <TrisaLogo />
                <Text fontWeight="bold" position="absolute" right="21px" top="72px">
                  Verified <br />
                  <span>
                    VASP{' '}
                    <BsCheckCircleFill
                      fontSize={20}
                      style={{ display: 'inline' }}
                      color="#34A853"
                    />
                  </span>
                </Text>
              </Box>
              <Box>
                <Button disabled>Download</Button>
              </Box>
            </Stack>
          </Box>
        </Stack>
      </Card.CardBody>
    </Card>
  );
};

export default TrisaVerifiedLogo;
