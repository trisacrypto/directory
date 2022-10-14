import { FiCompass } from 'react-icons/fi';
import { FaRegLightbulb, FaRegMoneyBillAlt } from 'react-icons/fa';
import { BiCertification } from 'react-icons/bi';
import { IconType } from 'react-icons';
import { CheckCircleIcon } from '@chakra-ui/icons';
import { ComponentWithAs, IconProps } from '@chakra-ui/react';
import { BsFillInfoCircleFill } from 'react-icons/bs';
import { t } from '@lingui/macro';
// import CertificateManagementIcon from 'assets/certificate-management.svg';

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
    title: t`Overview`,
    icon: FiCompass,
    activated: true,
    path: '/dashboard/overview'
  },
  {
    title: t`Certificate Management`,
    icon: BiCertification,
    activated: true,
    // path: '/dashboard/certificate-management',
    children: [
      {
        title: t`Certificate Management`,
        icon: FaRegMoneyBillAlt,
        path: '/dashboard/certificate-management',
        activated: true
      },
      {
        title: t`Certificate Registration`,
        icon: CheckCircleIcon,
        path: '/dashboard/certificate/registration',
        activated: true
      },
      {
        title: t`Certificate Details`,
        icon: BsFillInfoCircleFill,
        path: '/dashboard/certificate/details',
        activated: false
      }
    ]
  },
  {
    title: t`Technical Resources`,
    icon: FaRegLightbulb
  }
];

export default MENU;
