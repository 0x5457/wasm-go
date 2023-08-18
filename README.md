# wasm_go

```sh

curl -s https://registry-cdn.wapm.io/contents/liftm/cowsay/0.2.2/target/wasm32-wasi/release/cowsay.wasm | wasmgo

```

# tests
> Test cases from https://github.com/WebAssembly/testsuite
### run tests
```sh
cd tests
node --experimental-wasi-unstable-preview1 setup_suite.js

go test
```

- [X] address.wast
- [ ] align.wast
- [ ] binary-leb128.wast
- [ ] binary.wast
- [ ] block.wast
- [ ] br.wast
- [ ] br_if.wast
- [ ] br_table.wast
- [ ] bulk.wast
- [ ] call.wast
- [ ] call_indirect.wast
- [ ] comments.wast
- [ ] const.wast
- [ ] conversions.wast
- [ ] custom.wast
- [ ] data.wast
- [ ] elem.wast
- [ ] endianness.wast
- [ ] exports.wast
- [X] f32.wast
- [ ] f32_bitwise.wast
- [ ] f32_cmp.wast
- [X] f64.wast
- [ ] f64_bitwise.wast
- [ ] f64_cmp.wast
- [ ] fac.wast
- [ ] float_exprs.wast
- [ ] float_literals.wast
- [ ] float_memory.wast
- [ ] float_misc.wast
- [ ] forward.wast
- [ ] func.wast
- [ ] func_ptrs.wast
- [ ] global.wast
- [X] i32.wast
- [X] i64.wast
- [ ] if.wast
- [ ] imports.wast
- [ ] inline-module.wast
- [ ] int_exprs.wast
- [ ] int_literals.wast
- [ ] labels.wast
- [ ] left-to-right.wast
- [ ] linking.wast
- [ ] load.wast
- [ ] local_get.wast
- [ ] local_set.wast
- [ ] local_tee.wast
- [ ] loop.wast
- [ ] memory.wast
- [ ] memory_copy.wast
- [ ] memory_fill.wast
- [ ] memory_grow.wast
- [ ] memory_init.wast
- [ ] memory_redundancy.wast
- [ ] memory_size.wast
- [ ] memory_trap.wast
- [ ] names.wast
- [ ] nop.wast
- [ ] ref_func.wast
- [ ] ref_is_null.wast
- [ ] ref_null.wast
- [ ] return.wast
- [ ] select.wast
- [ ] simd_address.wast
- [ ] simd_align.wast
- [ ] simd_bit_shift.wast
- [ ] simd_bitwise.wast
- [ ] simd_boolean.wast
- [ ] simd_const.wast
- [ ] simd_conversions.wast
- [ ] simd_f32x4.wast
- [ ] simd_f32x4_arith.wast
- [ ] simd_f32x4_cmp.wast
- [ ] simd_f32x4_pmin_pmax.wast
- [ ] simd_f32x4_rounding.wast
- [ ] simd_f64x2.wast
- [ ] simd_f64x2_arith.wast
- [ ] simd_f64x2_cmp.wast
- [ ] simd_f64x2_pmin_pmax.wast
- [ ] simd_f64x2_rounding.wast
- [ ] simd_i16x8_arith.wast
- [ ] simd_i16x8_arith2.wast
- [ ] simd_i16x8_cmp.wast
- [ ] simd_i16x8_extadd_pairwise_i8x16.wast
- [ ] simd_i16x8_extmul_i8x16.wast
- [ ] simd_i16x8_q15mulr_sat_s.wast
- [ ] simd_i16x8_sat_arith.wast
- [ ] simd_i32x4_arith.wast
- [ ] simd_i32x4_arith2.wast
- [ ] simd_i32x4_cmp.wast
- [ ] simd_i32x4_dot_i16x8.wast
- [ ] simd_i32x4_extadd_pairwise_i16x8.wast
- [ ] simd_i32x4_extmul_i16x8.wast
- [ ] simd_i32x4_trunc_sat_f32x4.wast
- [ ] simd_i32x4_trunc_sat_f64x2.wast
- [ ] simd_i64x2_arith.wast
- [ ] simd_i64x2_arith2.wast
- [ ] simd_i64x2_cmp.wast
- [ ] simd_i64x2_extmul_i32x4.wast
- [ ] simd_i8x16_arith.wast
- [ ] simd_i8x16_arith2.wast
- [ ] simd_i8x16_cmp.wast
- [ ] simd_i8x16_sat_arith.wast
- [ ] simd_int_to_int_extend.wast
- [ ] simd_lane.wast
- [ ] simd_linking.wast
- [ ] simd_load.wast
- [ ] simd_load16_lane.wast
- [ ] simd_load32_lane.wast
- [ ] simd_load64_lane.wast
- [ ] simd_load8_lane.wast
- [ ] simd_load_extend.wast
- [ ] simd_load_splat.wast
- [ ] simd_load_zero.wast
- [ ] simd_splat.wast
- [ ] simd_store.wast
- [ ] simd_store16_lane.wast
- [ ] simd_store32_lane.wast
- [ ] simd_store64_lane.wast
- [ ] simd_store8_lane.wast
- [ ] skip-stack-guard-page.wast
- [ ] stack.wast
- [ ] start.wast
- [ ] store.wast
- [ ] switch.wast
- [ ] table-sub.wast
- [ ] table.wast
- [ ] table_copy.wast
- [ ] table_fill.wast
- [ ] table_get.wast
- [ ] table_grow.wast
- [ ] table_init.wast
- [ ] table_set.wast
- [ ] table_size.wast
- [ ] token.wast
- [ ] tokens.wast
- [ ] traps.wast
- [ ] type.wast
- [ ] unreachable.wast
- [ ] unreached-invalid.wast
- [ ] unreached-valid.wast
- [ ] unwind.wast
- [ ] utf8-custom-section-id.wast
- [ ] utf8-import-field.wast
- [ ] utf8-import-module.wast
- [ ] utf8-invalid-encoding.wast