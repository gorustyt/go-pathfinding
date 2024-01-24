// Code generated by "stringer -output binres_string.go -type ResType,DataType"; DO NOT EDIT.

package binres

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ResNull-0]
	_ = x[ResStringPool-1]
	_ = x[ResTable-2]
	_ = x[ResXML-3]
	_ = x[ResXMLStartNamespace-256]
	_ = x[ResXMLEndNamespace-257]
	_ = x[ResXMLStartElement-258]
	_ = x[ResXMLEndElement-259]
	_ = x[ResXMLCharData-260]
	_ = x[ResXMLResourceMap-384]
	_ = x[ResTablePackage-512]
	_ = x[ResTableType-513]
	_ = x[ResTableTypeSpec-514]
	_ = x[ResTableLibrary-515]
}

const (
	_ResType_name_0 = "ResNullResStringPoolResTableResXML"
	_ResType_name_1 = "ResXMLStartNamespaceResXMLEndNamespaceResXMLStartElementResXMLEndElementResXMLCharData"
	_ResType_name_2 = "ResXMLResourceMap"
	_ResType_name_3 = "ResTablePackageResTableTypeResTableTypeSpecResTableLibrary"
)

var (
	_ResType_index_0 = [...]uint8{0, 7, 20, 28, 34}
	_ResType_index_1 = [...]uint8{0, 20, 38, 56, 72, 86}
	_ResType_index_3 = [...]uint8{0, 15, 27, 43, 58}
)

func (i ResType) String() string {
	switch {
	case i <= 3:
		return _ResType_name_0[_ResType_index_0[i]:_ResType_index_0[i+1]]
	case 256 <= i && i <= 260:
		i -= 256
		return _ResType_name_1[_ResType_index_1[i]:_ResType_index_1[i+1]]
	case i == 384:
		return _ResType_name_2
	case 512 <= i && i <= 515:
		i -= 512
		return _ResType_name_3[_ResType_index_3[i]:_ResType_index_3[i+1]]
	default:
		return "ResType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DataNull-0]
	_ = x[DataReference-1]
	_ = x[DataAttribute-2]
	_ = x[DataString-3]
	_ = x[DataFloat-4]
	_ = x[DataDimension-5]
	_ = x[DataFraction-6]
	_ = x[DataDynamicReference-7]
	_ = x[DataIntDec-16]
	_ = x[DataIntHex-17]
	_ = x[DataIntBool-18]
	_ = x[DataIntColorARGB8-28]
	_ = x[DataIntColorRGB8-29]
	_ = x[DataIntColorARGB4-30]
	_ = x[DataIntColorRGB4-31]
}

const (
	_DataType_name_0 = "DataNullDataReferenceDataAttributeDataStringDataFloatDataDimensionDataFractionDataDynamicReference"
	_DataType_name_1 = "DataIntDecDataIntHexDataIntBool"
	_DataType_name_2 = "DataIntColorARGB8DataIntColorRGB8DataIntColorARGB4DataIntColorRGB4"
)

var (
	_DataType_index_0 = [...]uint8{0, 8, 21, 34, 44, 53, 66, 78, 98}
	_DataType_index_1 = [...]uint8{0, 10, 20, 31}
	_DataType_index_2 = [...]uint8{0, 17, 33, 50, 66}
)

func (i DataType) String() string {
	switch {
	case i <= 7:
		return _DataType_name_0[_DataType_index_0[i]:_DataType_index_0[i+1]]
	case 16 <= i && i <= 18:
		i -= 16
		return _DataType_name_1[_DataType_index_1[i]:_DataType_index_1[i+1]]
	case 28 <= i && i <= 31:
		i -= 28
		return _DataType_name_2[_DataType_index_2[i]:_DataType_index_2[i+1]]
	default:
		return "DataType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}