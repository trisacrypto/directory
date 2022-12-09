import { Alert, AlertDescription, AlertIcon, Button, chakra, HStack } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { useNavigate } from 'react-router-dom';

const CertificateValidityAlert = ({ hasValidMainnet, hasValidTestnet }: any) => {
  const { jumpToLastStep } = useCertificateStepper();
  const navigate = useNavigate();

  const handleJumpToLastStep = () => {
    navigate('/dashboard/certificate/registration');
    jumpToLastStep();
  };

  return (
    <>
      {hasValidMainnet && hasValidTestnet ? (
        <Alert bg="#D8EAF6" borderRadius={'10px'} mb={2}>
          <AlertIcon />
          <HStack justifyContent={'space-between'} w="100%">
            <AlertDescription>
              <Trans>
                The Organization has a current and valid{' '}
                <chakra.span fontWeight={700}>Mainnet</chakra.span> and{' '}
                <chakra.span fontWeight={700}>Tesnet</chakra.span> Identity Certificate.
              </Trans>
            </AlertDescription>

            <Button
              onClick={handleJumpToLastStep}
              border={'1px solid white'}
              width={142}
              px={8}
              as={'a'}
              borderRadius={0}
              color="#fff"
              cursor="pointer"
              bg="#000"
              _hover={{ bg: '#000000D1' }}>
              <Trans>View/Edit</Trans>
            </Button>
          </HStack>
        </Alert>
      ) : null}

      {hasValidMainnet ? (
        <Alert bg="#D8EAF6" borderRadius={'10px'} mb={2}>
          <AlertIcon />
          <HStack justifyContent={'space-between'} w="100%">
            <AlertDescription>
              <Trans>
                The Organization has a current and valid{' '}
                <chakra.span fontWeight={700}>Mainnet</chakra.span> Identity Certificate.
              </Trans>
            </AlertDescription>
            <Button
              onClick={handleJumpToLastStep}
              border={'1px solid white'}
              width={142}
              px={8}
              as={'a'}
              borderRadius={0}
              color="#fff"
              cursor="pointer"
              bg="#000"
              _hover={{ bg: '#000000D1' }}>
              <Trans>View/Edit</Trans>
            </Button>
          </HStack>
        </Alert>
      ) : null}

      {hasValidTestnet ? (
        <Alert bg="#D8EAF6" borderRadius={'10px'} mb={2}>
          <AlertIcon />
          <HStack justifyContent={'space-between'} w="100%">
            <AlertDescription>
              <Trans>
                The Organization has a current and valid{' '}
                <chakra.span fontWeight={700}>Testnet</chakra.span> Identity Certificate.
              </Trans>
            </AlertDescription>
            <Button
              onClick={handleJumpToLastStep}
              border={'1px solid white'}
              width={142}
              px={8}
              as={'a'}
              borderRadius={0}
              color="#fff"
              cursor="pointer"
              bg="#000"
              _hover={{ bg: '#000000D1' }}>
              <Trans>View/Edit</Trans>
            </Button>
          </HStack>
        </Alert>
      ) : null}
    </>
  );
};

export default CertificateValidityAlert;
