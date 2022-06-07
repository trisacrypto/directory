import { Flex, Heading, Stack, VStack, Text, Box } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import Footer from 'components/Footer/LandingFooter';
import LandingHeader from 'components/Header/LandingHeader';
import MembershipGuideCard from 'components/MembershipGuideCard';
import { NavBar } from 'components/Navbar/Landing/Nav';
import React from 'react';

const MembershipGuideText = [
  {
    stepNumber: 1,
    header: t`create your account`,
    description: t`Create your TRISA account with your VASP email address. Add collaborators in your organization.`,
    buttonText: t`Create Account`
  },
  {
    stepNumber: 2,
    header: t`complete VASP verification`,
    description: t`Complete the multi-part TRISA verification form and due diligence process. Once approved, gain access to the Testnet.`,
    buttonText: t`Learn More`
  },
  {
    stepNumber: 3,
    header: t`Integrate and Comply`,
    description: t`Set up your TRISA node or integrate with a 3rd-party Travel Rule solution. Complete testing and move to production.`,
    buttonText: t`Learn More`
  }
];

const MembershipGuide = () => {
  return (
    <Box minHeight="100vh">
      <LandingHeader />
      <section>
        <Stack>
          <Flex
            bgGradient="linear-gradient(90.17deg, rgba(35, 167, 224, 0.85) 3.85%, rgba(27, 206, 159, 0.55) 96.72%);"
            color="white"
            width="100%"
            minHeight={286}
            justifyContent="center"
            direction="column"
            paddingY={{ base: 12, md: 16 }}
            px="1rem"
            fontSize={'xl'}>
            <Stack textAlign={'center'} color="white" spacing={{ base: 3 }}>
              <VStack spacing={1}>
                <Heading
                  fontWeight={700}
                  fontFamily="Open Sans, sans-serif !important"
                  fontSize={{ md: '4xl', sm: '2xl' }}
                  color="#fff">
                  <Trans id="Welcome to TRISA’s network of Certified VASPs.">
                    Welcome to TRISA’s network of Certified VASPs.
                  </Trans>
                </Heading>
                <Text as="p" mt={2}>
                  <Trans id="Learn about the three-step certificaton process">
                    Learn about the three-step certificaton process
                  </Trans>
                </Text>
                <Text as="p" mt={2}>
                  <Trans id="and create your account today.">and create your account today.</Trans>
                </Text>
              </VStack>
            </Stack>
          </Flex>
          <Stack
            justifyContent={'center'}
            alignItems={['center', 'stretch']}
            spacing={10}
            direction={['column', 'row']}
            py={'2rem'}>
            {MembershipGuideText.map(({ stepNumber, header, description, buttonText }) => (
              <React.Fragment key={stepNumber}>
                <MembershipGuideCard
                  stepNumber={stepNumber}
                  header={header}
                  description={description}
                  buttonText={buttonText}
                />
              </React.Fragment>
            ))}
          </Stack>
        </Stack>
      </section>
      <Footer />
    </Box>
  );
};

export default MembershipGuide;
