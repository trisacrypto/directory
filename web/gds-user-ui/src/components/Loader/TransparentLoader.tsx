import { Spinner, Flex, Text, VStack } from '@chakra-ui/react';

interface LoaderProps {
  title?: any;
  opacity?: 'md' | 'full';
}
// This loader is for containers that are loading
const TransparentLoader: React.FC<LoaderProps> = (props) => {
  const { title, opacity } = props;
  const bgColor = opacity === 'full' ? 'rgba(255, 255, 255, 255)' : 'rgba(255,255,255,0.7)';
  return (
    <Flex
      position={'fixed'}
      width={'100%'}
      left={0}
      top={0}
      right={0}
      bottom={0}
      backgroundColor={bgColor}
      zIndex={1}>
      <VStack spacing={4} m={'auto'} py={10}>
        <Spinner color="blue.500" size="xl" />
        <Text>{title}</Text>
      </VStack>
    </Flex>
  );
};
TransparentLoader.defaultProps = {
  title: 'Session expired , reconnecting...'
};

export default TransparentLoader;
