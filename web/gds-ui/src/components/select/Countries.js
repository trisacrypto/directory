import 'i18n-iso-countries'
import React, { useContext } from 'react';
import LanguageContext from "../../contexts/LanguageContext";


var countries = require("i18n-iso-countries");

const Countries = () => {

  const context = useContext(LanguageContext);
  countries.registerLocale(require(`i18n-iso-countries/langs/${context.language}.json`));
  
  return (
    <>
    <option value=""></option>
    <option value="AF">{countries.getName("AF", context.language)}</option>
    <option value="AX">{countries.getName("AX", context.language)}</option>
    <option value="AL">{countries.getName("AL", context.language)}</option>
    <option value="DZ">{countries.getName("DZ", context.language)}</option>
    <option value="AS">{countries.getName("AS", context.language)}</option>
    <option value="AD">{countries.getName("AD", context.language)}</option>
    <option value="AO">{countries.getName("AO", context.language)}</option>
    <option value="AI">{countries.getName("AI", context.language)}</option>
    <option value="AQ">{countries.getName("AQ", context.language)}</option>
    <option value="AG">{countries.getName("AG", context.language)}</option>
    <option value="AR">{countries.getName("AR", context.language)}</option>
    <option value="AM">{countries.getName("AM", context.language)}</option>
    <option value="AW">{countries.getName("AW", context.language)}</option>
    <option value="AU">{countries.getName("AU", context.language)}</option>
    <option value="AT">{countries.getName("AT", context.language)}</option>
    <option value="AZ">{countries.getName("AZ", context.language)}</option>
    <option value="BS">{countries.getName("BS", context.language)}</option>
    <option value="BH">{countries.getName("BH", context.language)}</option>
    <option value="BD">{countries.getName("BD", context.language)}</option>
    <option value="BB">{countries.getName("BB", context.language)}</option>
    <option value="BY">{countries.getName("BY", context.language)}</option>
    <option value="BE">{countries.getName("BE", context.language)}</option>
    <option value="BZ">{countries.getName("BZ", context.language)}</option>
    <option value="BJ">{countries.getName("BJ", context.language)}</option>
    <option value="BM">{countries.getName("BM", context.language)}</option>
    <option value="BT">{countries.getName("BT", context.language)}</option>
    <option value="BO">{countries.getName("BO", context.language)}</option>
    <option value="BQ">{countries.getName("BQ", context.language)}</option>
    <option value="BA">{countries.getName("BA", context.language)}</option>
    <option value="BW">{countries.getName("BW", context.language)}</option>
    <option value="BV">{countries.getName("BV", context.language)}</option>
    <option value="BR">{countries.getName("BR", context.language)}</option>
    <option value="IO">{countries.getName("IO", context.language)}</option>
    <option value="BN">{countries.getName("BN", context.language)}</option>
    <option value="BG">{countries.getName("BG", context.language)}</option>
    <option value="BF">{countries.getName("BF", context.language)}</option>
    <option value="BI">{countries.getName("BI", context.language)}</option>
    <option value="CV">{countries.getName("CV", context.language)}</option>
    <option value="KH">{countries.getName("KH", context.language)}</option>
    <option value="CM">{countries.getName("CM", context.language)}</option>
    <option value="CA">{countries.getName("CA", context.language)}</option>
    <option value="KY">{countries.getName("KY", context.language)}</option>
    <option value="CF">{countries.getName("CF", context.language)}</option>
    <option value="TD">{countries.getName("TD", context.language)}</option>
    <option value="CL">{countries.getName("CL", context.language)}</option>
    <option value="CN">{countries.getName("CN", context.language)}</option>
    <option value="CX">{countries.getName("CX", context.language)}</option>
    <option value="CC">{countries.getName("CC", context.language)}</option>
    <option value="CO">{countries.getName("CO", context.language)}</option>
    <option value="KM">{countries.getName("KM", context.language)}</option>
    <option value="CD">{countries.getName("CD", context.language)}</option>
    <option value="CG">{countries.getName("CG", context.language)}</option>
    <option value="CK">{countries.getName("CK", context.language)}</option>
    <option value="CR">{countries.getName("CR", context.language)}</option>
    <option value="HR">{countries.getName("HR", context.language)}</option>
    <option value="CU">{countries.getName("CU", context.language)}</option>
    <option value="CW">{countries.getName("CW", context.language)}</option>
    <option value="CY">{countries.getName("CY", context.language)}</option>
    <option value="CZ">{countries.getName("CZ", context.language)}</option>
    <option value="CI">{countries.getName("CI", context.language)}</option>
    <option value="DK">{countries.getName("DK", context.language)}</option>
    <option value="DJ">{countries.getName("DJ", context.language)}</option>
    <option value="DM">{countries.getName("DM", context.language)}</option>
    <option value="DO">{countries.getName("DO", context.language)}</option>
    <option value="EC">{countries.getName("EC", context.language)}</option>
    <option value="EG">{countries.getName("EG", context.language)}</option>
    <option value="SV">{countries.getName("SV", context.language)}</option>
    <option value="GQ">{countries.getName("GQ", context.language)}</option>
    <option value="ER">{countries.getName("ER", context.language)}</option>
    <option value="EE">{countries.getName("EE", context.language)}</option>
    <option value="SZ">{countries.getName("SZ", context.language)}</option>
    <option value="ET">{countries.getName("ET", context.language)}</option>
    <option value="FK">{countries.getName("FK", context.language)}</option>
    <option value="FO">{countries.getName("FO", context.language)}</option>
    <option value="FJ">{countries.getName("FJ", context.language)}</option>
    <option value="FI">{countries.getName("FI", context.language)}</option>
    <option value="FR">{countries.getName("FR", context.language)}</option>
    <option value="GF">{countries.getName("GF", context.language)}</option>
    <option value="PF">{countries.getName("PF", context.language)}</option>
    <option value="TF">{countries.getName("TF", context.language)}</option>
    <option value="GA">{countries.getName("GA", context.language)}</option>
    <option value="GM">{countries.getName("GM", context.language)}</option>
    <option value="GE">{countries.getName("GE", context.language)}</option>
    <option value="DE">{countries.getName("DE", context.language)}</option>
    <option value="GH">{countries.getName("GH", context.language)}</option>
    <option value="GI">{countries.getName("GI", context.language)}</option>
    <option value="GR">{countries.getName("GR", context.language)}</option>
    <option value="GL">{countries.getName("GL", context.language)}</option>
    <option value="GD">{countries.getName("GD", context.language)}</option>
    <option value="GP">{countries.getName("GP", context.language)}</option>
    <option value="GU">{countries.getName("GU", context.language)}</option>
    <option value="GT">{countries.getName("GT", context.language)}</option>
    <option value="GG">{countries.getName("GG", context.language)}</option>
    <option value="GN">{countries.getName("GN", context.language)}</option>
    <option value="GW">{countries.getName("GW", context.language)}</option>
    <option value="GY">{countries.getName("GY", context.language)}</option>
    <option value="HT">{countries.getName("HT", context.language)}</option>
    <option value="HM">{countries.getName("HM", context.language)}</option>
    <option value="VA">{countries.getName("VA", context.language)}</option>
    <option value="HN">{countries.getName("HN", context.language)}</option>
    <option value="HK">{countries.getName("HK", context.language)}</option>
    <option value="HU">{countries.getName("HU", context.language)}</option>
    <option value="IS">{countries.getName("IS", context.language)}</option>
    <option value="IN">{countries.getName("IN", context.language)}</option>
    <option value="ID">{countries.getName("ID", context.language)}</option>
    <option value="IR">{countries.getName("IR", context.language)}</option>
    <option value="IQ">{countries.getName("IQ", context.language)}</option>
    <option value="IE">{countries.getName("IE", context.language)}</option>
    <option value="IM">{countries.getName("IM", context.language)}</option>
    <option value="IL">{countries.getName("IL", context.language)}</option>
    <option value="IT">{countries.getName("IT", context.language)}</option>
    <option value="JM">{countries.getName("JM", context.language)}</option>
    <option value="JP">{countries.getName("JP", context.language)}</option>
    <option value="JE">{countries.getName("JE", context.language)}</option>
    <option value="JO">{countries.getName("JO", context.language)}</option>
    <option value="KZ">{countries.getName("KZ", context.language)}</option>
    <option value="KE">{countries.getName("KE", context.language)}</option>
    <option value="KI">{countries.getName("KI", context.language)}</option>
    <option value="KP">{countries.getName("KP", context.language)}</option>
    <option value="KR">{countries.getName("KR", context.language)}</option>
    <option value="KW">{countries.getName("KW", context.language)}</option>
    <option value="KG">{countries.getName("KG", context.language)}</option>
    <option value="LA">{countries.getName("LA", context.language)}</option>
    <option value="LV">{countries.getName("LV", context.language)}</option>
    <option value="LB">{countries.getName("LB", context.language)}</option>
    <option value="LS">{countries.getName("LS", context.language)}</option>
    <option value="LR">{countries.getName("LR", context.language)}</option>
    <option value="LY">{countries.getName("LY", context.language)}</option>
    <option value="LI">{countries.getName("LI", context.language)}</option>
    <option value="LT">{countries.getName("LT", context.language)}</option>
    <option value="LU">{countries.getName("LU", context.language)}</option>
    <option value="MO">{countries.getName("MO", context.language)}</option>
    <option value="MG">{countries.getName("MG", context.language)}</option>
    <option value="MW">{countries.getName("MW", context.language)}</option>
    <option value="MY">{countries.getName("MY", context.language)}</option>
    <option value="MV">{countries.getName("MV", context.language)}</option>
    <option value="ML">{countries.getName("ML", context.language)}</option>
    <option value="MT">{countries.getName("MT", context.language)}</option>
    <option value="MH">{countries.getName("MH", context.language)}</option>
    <option value="MQ">{countries.getName("MQ", context.language)}</option>
    <option value="MR">{countries.getName("MR", context.language)}</option>
    <option value="MU">{countries.getName("MU", context.language)}</option>
    <option value="YT">{countries.getName("YT", context.language)}</option>
    <option value="MX">{countries.getName("MX", context.language)}</option>
    <option value="FM">{countries.getName("FM", context.language)}</option>
    <option value="MD">{countries.getName("MD", context.language)}</option>
    <option value="MC">{countries.getName("MC", context.language)}</option>
    <option value="MN">{countries.getName("MN", context.language)}</option>
    <option value="ME">{countries.getName("ME", context.language)}</option>
    <option value="MS">{countries.getName("MS", context.language)}</option>
    <option value="MA">{countries.getName("MA", context.language)}</option>
    <option value="MZ">{countries.getName("MZ", context.language)}</option>
    <option value="MM">{countries.getName("MM", context.language)}</option>
    <option value="NA">{countries.getName("NA", context.language)}</option>
    <option value="NR">{countries.getName("NR", context.language)}</option>
    <option value="NP">{countries.getName("NP", context.language)}</option>
    <option value="NL">{countries.getName("NL", context.language)}</option>
    <option value="NC">{countries.getName("NC", context.language)}</option>
    <option value="NZ">{countries.getName("NZ", context.language)}</option>
    <option value="NI">{countries.getName("NI", context.language)}</option>
    <option value="NE">{countries.getName("NE", context.language)}</option>
    <option value="NG">{countries.getName("NG", context.language)}</option>
    <option value="NU">{countries.getName("NU", context.language)}</option>
    <option value="NF">{countries.getName("NF", context.language)}</option>
    <option value="MP">{countries.getName("MP", context.language)}</option>
    <option value="NO">{countries.getName("NO", context.language)}</option>
    <option value="OM">{countries.getName("OM", context.language)}</option>
    <option value="PK">{countries.getName("PK", context.language)}</option>
    <option value="PW">{countries.getName("PW", context.language)}</option>
    <option value="PS">{countries.getName("PS", context.language)}</option>
    <option value="PA">{countries.getName("PA", context.language)}</option>
    <option value="PG">{countries.getName("PG", context.language)}</option>
    <option value="PY">{countries.getName("PY", context.language)}</option>
    <option value="PE">{countries.getName("PE", context.language)}</option>
    <option value="PH">{countries.getName("PH", context.language)}</option>
    <option value="PN">{countries.getName("PN", context.language)}</option>
    <option value="PL">{countries.getName("PL", context.language)}</option>
    <option value="PT">{countries.getName("PT", context.language)}</option>
    <option value="PR">{countries.getName("PR", context.language)}</option>
    <option value="QA">{countries.getName("QA", context.language)}</option>
    <option value="MK">{countries.getName("MK", context.language)}</option>
    <option value="RO">{countries.getName("RO", context.language)}</option>
    <option value="RU">{countries.getName("RU", context.language)}</option>
    <option value="RW">{countries.getName("RW", context.language)}</option>
    <option value="RE">{countries.getName("RE", context.language)}</option>
    <option value="BL">{countries.getName("BL", context.language)}</option>
    <option value="SH">{countries.getName("SH", context.language)}</option>
    <option value="KN">{countries.getName("KN", context.language)}</option>
    <option value="LC">{countries.getName("LC", context.language)}</option>
    <option value="MF">{countries.getName("MF", context.language)}</option>
    <option value="PM">{countries.getName("PM", context.language)}</option>
    <option value="VC">{countries.getName("VC", context.language)}</option>
    <option value="WS">{countries.getName("WS", context.language)}</option>
    <option value="SM">{countries.getName("SM", context.language)}</option>
    <option value="ST">{countries.getName("ST", context.language)}</option>
    <option value="SA">{countries.getName("SA", context.language)}</option>
    <option value="SN">{countries.getName("SN", context.language)}</option>
    <option value="RS">{countries.getName("RS", context.language)}</option>
    <option value="SC">{countries.getName("SC", context.language)}</option>
    <option value="SL">{countries.getName("SL", context.language)}</option>
    <option value="SG">{countries.getName("SG", context.language)}</option>
    <option value="SX">{countries.getName("SX", context.language)}</option>
    <option value="SK">{countries.getName("SK", context.language)}</option>
    <option value="SI">{countries.getName("SI", context.language)}</option>
    <option value="SB">{countries.getName("SB", context.language)}</option>
    <option value="SO">{countries.getName("SO", context.language)}</option>
    <option value="ZA">{countries.getName("ZA", context.language)}</option>
    <option value="GS">{countries.getName("GS", context.language)}</option>
    <option value="SS">{countries.getName("SS", context.language)}</option>
    <option value="ES">{countries.getName("ES", context.language)}</option>
    <option value="LK">{countries.getName("LK", context.language)}</option>
    <option value="SD">{countries.getName("SD", context.language)}</option>
    <option value="SR">{countries.getName("SR", context.language)}</option>
    <option value="SJ">{countries.getName("SJ", context.language)}</option>
    <option value="SE">{countries.getName("SE", context.language)}</option>
    <option value="CH">{countries.getName("CH", context.language)}</option>
    <option value="SY">{countries.getName("SY", context.language)}</option>
    <option value="TW">{countries.getName("TW", context.language)}</option>
    <option value="TJ">{countries.getName("TJ", context.language)}</option>
    <option value="TZ">{countries.getName("TZ", context.language)}</option>
    <option value="TH">{countries.getName("TH", context.language)}</option>
    <option value="TL">{countries.getName("TL", context.language)}</option>
    <option value="TG">{countries.getName("TG", context.language)}</option>
    <option value="TK">{countries.getName("TK", context.language)}</option>
    <option value="TO">{countries.getName("TO", context.language)}</option>
    <option value="TT">{countries.getName("TT", context.language)}</option>
    <option value="TN">{countries.getName("TN", context.language)}</option>
    <option value="TR">{countries.getName("TR", context.language)}</option>
    <option value="TM">{countries.getName("TM", context.language)}</option>
    <option value="TC">{countries.getName("TC", context.language)}</option>
    <option value="TV">{countries.getName("TV", context.language)}</option>
    <option value="UG">{countries.getName("UG", context.language)}</option>
    <option value="UA">{countries.getName("UA", context.language)}</option>
    <option value="AE">{countries.getName("AE", context.language)}</option>
    <option value="GB">{countries.getName("GB", context.language)}</option>
    <option value="US">{countries.getName("US", context.language)}</option>
    <option value="UY">{countries.getName("UY", context.language)}</option>
    <option value="UZ">{countries.getName("UZ", context.language)}</option>
    <option value="VU">{countries.getName("VU", context.language)}</option>
    <option value="VE">{countries.getName("VE", context.language)}</option>
    <option value="VN">{countries.getName("VN", context.language)}</option>
    <option value="VG">{countries.getName("VG", context.language)}</option>
    <option value="UM">{countries.getName("UM", context.language)}</option>
    <option value="VI">{countries.getName("VI", context.language)}</option>
    <option value="WF">{countries.getName("WF", context.language)}</option>
    <option value="EH">{countries.getName("EH", context.language)}</option>
    <option value="YE">{countries.getName("YE", context.language)}</option>
    <option value="ZM">{countries.getName("ZM", context.language)}</option>
    <option value="ZW">{countries.getName("ZW", context.language)}</option>
    </>
  )
}

export default Countries;