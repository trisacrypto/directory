// transparent background
import React from 'react';
import { Flex } from '@chakra-ui/react';
interface TransparentBackgroundProps {
  children: React.ReactNode;
  opacity?: 'md' | 'full';
}
export const TransparentBackground = ({ children }: TransparentBackgroundProps) => {
  // const bgColor = opacity === 'full' ? 'rgba(255, 255, 255, 255)' : 'rgba(255,255,255,0.7)';
  return (
    <Flex
      position="absolute"
      top="0"
      left="0"
      w="100%"
      h="100%"
      bg={'rgba(255,255,255,255)'}
      zIndex="8888">
      {children}
    </Flex>
  );
};
