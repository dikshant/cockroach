// Code generated by "stringer"; DO NOT EDIT.

package screl

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DescID-1]
	_ = x[IndexID-2]
	_ = x[ColumnFamilyID-3]
	_ = x[ColumnID-4]
	_ = x[ConstraintID-5]
	_ = x[Name-6]
	_ = x[ReferencedDescID-7]
	_ = x[Comment-8]
	_ = x[TemporaryIndexID-9]
	_ = x[SourceIndexID-10]
	_ = x[TargetStatus-11]
	_ = x[CurrentStatus-12]
	_ = x[Element-13]
	_ = x[Target-14]
	_ = x[ReferencedTypeIDs-15]
	_ = x[ReferencedSequenceIDs-16]
	_ = x[ReferencedFunctionIDs-17]
}

const _Attr_name = "DescIDIndexIDColumnFamilyIDColumnIDConstraintIDNameReferencedDescIDCommentTemporaryIndexIDSourceIndexIDTargetStatusCurrentStatusElementTargetReferencedTypeIDsReferencedSequenceIDsReferencedFunctionIDs"

var _Attr_index = [...]uint8{0, 6, 13, 27, 35, 47, 51, 67, 74, 90, 103, 115, 128, 135, 141, 158, 179, 200}

func (i Attr) String() string {
	i -= 1
	if i < 0 || i >= Attr(len(_Attr_index)-1) {
		return "Attr(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Attr_name[_Attr_index[i]:_Attr_index[i+1]]
}
