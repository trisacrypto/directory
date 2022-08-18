import { FiCompass } from 'react-icons/fi';
import { FaRegLightbulb, FaBook } from 'react-icons/fa';
import { HiOutlineUserGroup } from 'react-icons/hi';
import { BiCertification } from 'react-icons/bi';
import { IconType } from 'react-icons';
import { CheckCircleIcon } from '@chakra-ui/icons';
import { ComponentWithAs, IconProps } from '@chakra-ui/react';
import { BsFillInfoCircleFill, BsInfoCircle } from 'react-icons/bs';
import { MdContactSupport } from 'react-icons/md';
import { IoLogoSlack } from 'react-icons/io5';

type Menu = {
  title: string;
  icon?: IconType | ComponentWithAs<'svg', IconProps>;
  activated?: boolean;
  path?: string;
  children?: Menu[];
  isExternalLink?: boolean;
};

const MENU: Menu[] = [
  {
    title: 'Overview',
    icon: FiCompass,
    activated: true,
    path: '/dashboard/overview'
  },
  {
    title: 'Certificate Management',
    icon: BiCertification,
    activated: true,
    path: '/dashboard/certificate-management',
    children: [
      {
        title: 'Certificate Registration',
        icon: CheckCircleIcon,
        path: '/dashboard/certificate/registration',
        activated: true
      },
      {
        title: 'Certificate Details',
        icon: BsFillInfoCircleFill,
        path: '/dashboard/certificate/details',
        activated: false
      }
    ]
  },
  {
    title: 'Technical Resources',
    icon: FaRegLightbulb
  }
];

export default MENU;
