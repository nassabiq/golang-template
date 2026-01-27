#!/usr/bin/env bash
set -e

MODULE=$1
BASE=internal/modules/$MODULE

if [ -z "$MODULE" ]; then
    echo "❌ Error: Module name is required"
    echo "👉 Usage: $0 <module_name>"
    exit 1
fi

# Validate module name (should be lowercase alphanumeric)
if [[ ! "$MODULE" =~ ^[a-z][a-z0-9]*$ ]]; then
    echo "❌ Error: Module name must be lowercase alphanumeric and start with a letter"
    exit 1
fi

echo "🚀 Generating module: $MODULE"

# Create directory structure
mkdir -p \
  $BASE/domain \
  $BASE/usecase \
  $BASE/repository \
  $BASE/handler \
  $BASE/dto

# Helper function for template rendering with multiple substitutions
render () {
    local input=$1
    local output=$2
    
    # Create temporary file for processing
    local temp_file=$(mktemp)
    
    # Apply substitutions
    sed -e "s|github.com/nassabiq/golang-template/internal/modules/{{MODULE}}/|github.com/nassabiq/golang-template/internal/modules/$MODULE/|g" \
        -e "s|github.com/nassabiq/golang-template/proto/{{MODULE}}|github.com/nassabiq/golang-template/proto/$MODULE|g" \
        -e "s/{{MODULE}}/$(echo $MODULE | awk '{print toupper(substr($0,1,1)) substr($0,2)}')/g" \
        -e "s/{{MODULE|lower}}/$MODULE/g" \
        "$input" > "$temp_file"
    
    # Move to final destination
    mv "$temp_file" "$output"
}

# Generate domain layer files
render templates/domain_entity.go.tpl      $BASE/domain/entity.go
render templates/domain_repository.go.tpl  $BASE/domain/repository.go

# Generate usecase layer
render templates/usecase.go.tpl             $BASE/usecase/${MODULE}_usecase.go

# Generate repository layer
render templates/repository.go.tpl          $BASE/repository/${MODULE}_repository.go

# Generate handler layer
render templates/handler.go.tpl             $BASE/handler/${MODULE}_handler.go

# Generate DTO files
render templates/dto/request.go.tpl         $BASE/dto/request.go
render templates/dto/response.go.tpl        $BASE/dto/response.go

# Generate proto file
mkdir -p proto/$MODULE
render templates/proto.tpl                  proto/$MODULE/$MODULE.proto

echo "✅ Module '$MODULE' generated successfully!"
echo "📁 Generated files:"
echo "   ├── $BASE/domain/entity.go"
echo "   ├── $BASE/domain/repository.go"
echo "   ├── $BASE/usecase/${MODULE}_usecase.go"
echo "   ├── $BASE/repository/${MODULE}_repository.go"
echo "   ├── $BASE/handler/${MODULE}_handler.go"
echo "   ├── $BASE/dto/request.go"
echo "   ├── $BASE/dto/response.go"
echo "   └── proto/$MODULE/$MODULE.proto"
echo ""
echo "🗃️  Don't forget to:"
echo "   1. Create the database migration: make migrate-create name=create_${MODULE}s_table"
echo "   2. Generate proto: buf generate"
echo "   3. Register the module in your server/main.go"
