import * as React from 'react';
import { Stack, Text, VStack, Flex, Button } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

const GettingStartedSection = () => {
  return (
    <Flex justifyItems={'center'} mx={{ lg: 70, md: 40, sm: 10 }} py={10}>
      <Stack direction={['column', 'row']} spacing={5} fontSize={{ lg: 24, sm: 18, md: 20 }}>
        <VStack textAlign={'center'} bg={'#E5EDF1'} p={5}>
          <Text>
            <Trans id="Step 1">Step 1</Trans>
          </Text>
          <Text pt={3} fontWeight={'bold'}>
            <Trans id="Create account">Create account</Trans>
          </Text>
          <Text pb={3}>
            <Trans
              id=" create your TRISA account with your VASP email address. Add collaborators in your
            organization.">
              create your TRISA account with your VASP email address. Add collaborators in your
              organization.
            </Trans>
          </Text>
          <Button bg={'white'} color={'#221F1F'}>
            <Trans id="create account">create account</Trans>
          </Button>
        </VStack>
        <VStack textAlign={'center'} bg={'#E5EDF1'} p={5}>
          <Text> Step 1</Text>
          <Text pt={3} fontWeight={'bold'}>
            Complete VASP Verification
          </Text>
          <Text pb={3}>
            <Trans
              id="Complete the multi-part TRISA verification form and due diligence process. Once
            approved, gain access to the Testnet.">
              Complete the multi-part TRISA verification form and due diligence process. Once
              approved, gain access to the Testnet.
            </Trans>
          </Text>
          <Button bg={'white'} color={'#221F1F'}>
            <Trans id="Learn More">Learn More</Trans>
          </Button>
        </VStack>
        <VStack textAlign={'center'} bg={'#E5EDF1'} p={5}>
          <Text>
            <Trans id="Step 1">Step 1</Trans>
          </Text>
          <Text pt={3} fontWeight={'bold'}>
            <Trans id="Integrate and Comply">Integrate and Comply</Trans>
          </Text>
          <Text pb={3}>
            <Trans
              id="Set up your TRISA node or integrate with a 3rd-party Travel Rule solution. Complete
            testing and move to production.">
              Set up your TRISA node or integrate with a 3rd-party Travel Rule solution. Complete
              testing and move to production.
            </Trans>
          </Text>
          <Button bg={'white'} color={'#221F1F'}>
            <Trans id="Learn More">Learn More</Trans>
          </Button>
        </VStack>
      </Stack>
    </Flex>
  );
};

export default GettingStartedSection;
