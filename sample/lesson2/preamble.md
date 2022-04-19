{% macro commentary(now) %}

Our solution for $x$ was {{ now.x }} and our solution for $y$ depends on time. Currently, $t={{ now.time}}$ which means that $y={{ now.y }}$.

{% endmacro %}
